var raiki = new (function() {
  var connected = false;
  var connecting = false;
  var socket;
  var authHandler = function() {}
  var disconnectHandler = function() {}
  var _id;
  var actionCounter = 0;

  var actions = {};
  var events = {};

  this.id = function() {
    return _id;
  }

  this.connect = function(url, token) {
    if (connected || connecting) return;
    connecting = true;
    socket = new WebSocket(url);
    socket.onopen = function() {
      socket.send(JSON.stringify({auth: token}));
    }

    socket.onmessage = function(msg) {
      if (msg.data === 'h') return socket.send('h');

      var data;
      try {
        data = JSON.parse(msg.data);
      } catch (e) { console.error(e); return; }

      if (data.auth) {
        _id = data.auth;
        connected = true;
        connecting = false;
        authHandler(_id);
        return;
      }

      if (data.i) {
        return actions[data.i] && actions[data.i](data.d, data.m), void 0;
      } else if (data.e) {
        if (!events[data.e]) return;
        for (var i = 0; i < events[data.e].length; i++) {
          events[data.e][i](data.d);
        }
      }
    }

    socket.onclose = socket.onerror = disconnectHandler;
  }

  this.setDisconnectHandler = function(handler) {
    if (typeof handler !== 'function') return;
    disconnectHandler = handler;
    if (socket) {
      socket.onclose = socket.onerror = disconnectHandler;
    }
  }

  this.setAuthHandler = function(handler) {
    if (typeof handler !== 'function') return;
    authHandler = handler;
  }

  this.action = function(action, data, callback) {
    if (!connected || connecting) return;
    if (typeof action !== 'string' || typeof callback !== 'function') return;
    var c = ++actionCounter;
    actions[c] = callback;
    socket.send(JSON.stringify({
      i: c,
      t: action,
      d: data
    }));
  }

  this.on = function(event, callback) {
    if (typeof event !== 'string' || typeof callback !== 'function') return;
    events[event] = events[event] || [];
    events[event].push(callback);
  }
});

var s;
var app = angular.module('phynix', [angularDragula(angular)]);
app.controller('main', function($scope, dragulaService) {
  // Drgular
  dragulaService.options($scope, 'playlist', {
    accepts: function (el, container, handle) {
      return true;
    },
    copy: function(el, source) {
      // return !$(source).hasClass('items');
      return true;
    },
    invalid: function(el, handle) {
      return $(el).parent().hasClass('playlist');
    },
    revertOnSpill: true
  });

  $scope.$on('playlist.over', function(e, el, container) {
    container.addClass('drag-over');
  })
  $scope.$on('playlist.out', function(e, el, container) {
    container.removeClass('drag-over');
  });

  $scope.$on('playlist.drop', function(e, el, target, source, sibling) {
    console.info(arguments);

    if ($(target).hasClass('playlist')) {
      var pid = +$(target).attr('data-id');
      var cid = $(el).attr('data-cid');
      var ct = +$(el).attr('data-ct');
      console.info(pid, cid, ct);
      raiki.action('playlist.insert', { id: pid, contentID: cid, type: ct }, function(data, err) {
        console.info(data, err);
        if (err) return;
        var playlist = $scope.getPlaylist(pid);
        if (playlist) {
          playlist.items = playlist.items || [];
          playlist.items.unshift(data);
        }
        $scope.$apply();
      });
    } else if ($(target).hasClass('items')) {
      console.info(e, el, target, source, sibling);
    }
  });

  s = $scope;
  $scope.tabIndex = 0;
  $scope.views = {
    playlist: false,
    playlistViewType: 0
  }

  $scope.community = {
    chat: [],
    users: [],
    waitlist: [],
    media: {},
    meta: {}
  }

  $scope.user = {}

  $scope.users = {}

  $scope.sizePlaylistItems = function() {
    setTimeout(function() {
    $('#dynamic-styles').remove();
    var sheet = $('<style id="dynamic-styles"></style>').appendTo('head')[0];
    // Playlist grid resizing
    // var ratio = 16 / 9;
    // var width = $('.playlist-modal .right').prop('scrollWidth');
    // var min = 320;
    // var items = ~~(width / min);
    // var newWidth = width / items;
    // var newHeight = newWidth / ratio;
    // sheet.appendChild(document.createTextNode(`.playlist-modal .right.grid .item {
    //   width: ${newWidth}px!important;
    //   height: ${newHeight}px!important;
    // }`));

    var payload = '';
    var areas = {
      '.right .items': '.right.grid .items',
      '.right .results': '.right.grid .results',
      '.content .tray .history': '.content .tray .history.grid'
    };
    for (var i in areas) {
      var vars = calculate($(i).prop('scrollWidth'));
      payload += `${areas[i]} .item {
        width: ${vars.width}px!important;
        height: ${vars.height}px!important;
      }`;
    }

    function calculate(width) {
      console.info(width);
      var ratio = 16 / 9;
      var min = 320;
      var items = ~~(width / min);
      return {
        width: width / items,
        height: (width / items) / ratio
      };
    }

    sheet.appendChild(document.createTextNode(payload));
    }.bind(window), 1);
  }
  $(window).resize((function() { $scope.sizePlaylistItems(); }));
  $scope.sizePlaylistItems();
  $scope.$watch('views.playlistViewType', $scope.sizePlaylistItems);
  // $('.playlist-modal .right').resize($scope.sizePlaylistItems);
  // $(window).resize($scope.sizePlaylistItems);

  $scope.searchYoutube = function(query) {
    raiki.action('media.search', { query: query, type: 1 }, function(data, err) {
      if (err) return;
      $scope.views.searchResults = data;
      $scope.$apply();
    });
  }

  $scope.getEmbedURL = function(item) {
    if (!$scope.community.media && !item) return;
    // var media = item || $scope.community.media;
    // var id = $scope.community.media.item.contentID;

    var media, item, id;
    if (item) {
      media = item;
      id = item.contentID;
    } else if ($scope.community.media) {
      media = $scope.community.media
      item = media.item;
      id = item.contentID;
    }
    switch (item.type) {
      case 1:
        return `https://youtube.com/embed/${id}?autoplay=true&controls=2&start=${~~(media.elapsed)}`;
      case 2:
        return `https://w.soundcloud.com/player/?visual=true&url=http%3A%2F%2Fapi.soundcloud.com%2Ftracks%2F${id}&show_artwork=true&auto_play=true`;
      default:
        return '';
    }
  }

  $scope.joinWaitlist = function() {
    if ($('.media .user .waitlist').hasClass('loading')) return;
    $('.media .user .waitlist').addClass('loading');
    if ($('.media .user .waitlist').hasClass('active')) {
      return raiki.action('community.waitlist.leave', {}, function(data, err) {
        if (err) return;
        $('.media .user .waitlist').removeClass('loading');
        $('.media .user .waitlist').removeClass('active');
        $('.media .user .waitlist .text').text('Join Wait List');
      });
    }
    raiki.action('community.waitlist.join', {}, function(data, err) {
    $('.media .user .waitlist').removeClass('loading');
      if (err) return;
      $('.media .user .waitlist').addClass('active');
      $('.media .user .waitlist .text').text('Leave Wait List');
    });
  }

  $scope.toDuration = function(duration) {
    var minutes = ~~(duration / 60);
    var seconds = duration % 60;
    return (minutes < 10 ? '0' + minutes : minutes) + ':' + (seconds < 10 ? '0' + seconds : seconds);
  }

  $scope.getHost = function() {
    var staff = $scope.community.meta.staff;
    for (var i = 0; i < staff.length; i++) {
      if (staff[i].role === 6) return $scope.getUser(staff[i].userID);
    }
    return null;
  }

  $scope.getUser = function(id) {
    var users = $scope.community.users;
    for (var i = 0; i < users.length; i++) {
      if (users[i].id === id) return users[i];
    }
    return null;
  }

  $scope.getPlaylist = function(id) {
    if (!$scope.user.playlists) return;
    var playlists = $scope.user.playlists;
    for (var i = 0; i < playlists.length; i++) {
      if (playlists[i].id === id) return playlists[i];
    }
    return null;
  }

  $scope.getActivePlaylist = function() {
    if (!$scope.user.playlists) return;
    var playlists = $scope.user.playlists;
    for (var i = 0; i < playlists.length; i++) {
      if (playlists[i].active) return playlists[i];
    }
    return null;
  }

  $scope.iconActivatePlaylist = function(e) {
    var el = e.target;
    var parent = $(el).parent()
    var id = +$(el).parent().attr('data-id');
    var playlist = $scope.getActivePlaylist() || {};
    if (playlist.id === id || !parent.hasClass('viewing')) return;
    raiki.action('playlist.activate', { id: id }, function(data, err) {
      if (err) return;
      $scope.getActivePlaylist().active = false;
      $scope.getPlaylist(id).active = true;
      $scope.$apply();
    });
  }

  $scope.woot = function() {
    var votes = $('.media .votes');
    if (votes.hasClass('loading')) return;
    votes.addClass('loading');
    raiki.action('community.vote.woot', {}, function(data, err) {
      votes.removeClass('loading');
    });
  }

  $scope.save = function() {

  }

  $scope.meh = function() {
    var votes = $('.media .votes');
    if (votes.hasClass('loading')) return;
    votes.addClass('loading');
    raiki.action('community.vote.meh', {}, function(data, err) {
      votes.removeClass('loading');
    });
  }

  $scope.itemPreview = function(item) {
    var url = $scope.getEmbedURL(item);
    var modal = $scope.compileModal(item.artist + ' - ' + item.title, 'item-preview', `
      <iframe class="frame" src="${url}"></iframe>
    `, 'Close', false);

    modal.find('.footer .submit').click(function () {
      modal.fadeOut(100, function() {
        this.remove();
      })
    });

    modal.appendTo('body').fadeIn(100, function() {
      modal.find('input').focus();
    });
  }

  $scope.showPlaylistCreateModal = function() {
    var modal = $scope.compileModal('Create Playlist', 'create-playlist', `
      <form class="playlist-create-form">
        <input class="name-input" maxlength="30" placeholder="Name your playlist">
      </form>
    `, 'Create Playlist', true);

    modal.find('.playlist-create-form').submit(function(e) {
      e.preventDefault();
      if (modal.attr('loading') === 'true') return;
      modal
        .attr('loading', true)
        .find('.message').removeClass('showing');
      var name = modal.find('.name-input').val();
      raiki.action('playlist.create', { name: name }, function(data, err) {
        modal.attr('loading', false);
        if (err) return modal.find('.message').text(err).addClass('showing');
        console.info(data, err);
        modal.find('.header .mdi-close').click();

        $scope.user.playlists = $scope.user.playlists || [];
        $scope.user.playlists.map(function(playlist) {
          playlist.active = false;
        });
        $scope.user.playlists.push(data);
        $scope.views.playlistView = data.id;
        $scope.$apply();
      });
    });

    modal.find('.footer .submit').click(function() {
      modal.find('.playlist-create-form').submit();
    });

    modal.appendTo('body').fadeIn(100, function() {
      modal.find('input').focus();
    });
  }

  $scope.compileModal = function(title, css, content, submitText, cancellable) {
    var modal = $(`
      <div class="modal ${css}" cancellable="${!!cancellable}">
        <div class="container">
          <div class="header">
            <div class="text">${title}</div>
            <i class="mdi mdi-close"></i>
          </div>
          <div class="message"></div>
          <div class="loader"></div>
          <div class="content">
            ${content}
          </div>
          <div class="footer">
            <div class="button cancel">Cancel</div>
            <div class="button submit">${submitText}</div>
          </div>
        </div>
      </div>
    `).click(function(e) {
      if (e.target !== this) return;
      modal.fadeOut(100, function() {
        this.remove();
      })
    });

    modal.find('.footer .cancel, .header .mdi-close').click(function() {
      modal.fadeOut(100, function() {
        this.remove();
      })
    });

    return modal;
  }

  raiki.setDisconnectHandler(function(e) {
    var modal;
    switch (e.reason) {
    case 'session replaced':
      modal = $scope.compileModal('Session replaced', 'socket-replace', `
        <div class="text">Your session has been replaced.<br>This usually happens when you log in from another location. Click OK to refresh.</div>
      `, 'OK', false);
      break;
    default:
      modal = $scope.compileModal('Socket disconnected', 'socket-disconnect', `
        <div class="text">Socket disconnected.<br>Click OK to refresh</div>
      `, 'OK', false);
      break;
    }

    modal.click(function(e) {
      if (e.target !== this) return;
      window.location.reload();
    });

    modal.find('.footer .submit, .header .mdi-close').click(function() {
      window.location.reload();
    });

    modal.appendTo('body').fadeIn(100, function() {
      modal.find('input').focus();
    });
  });

  raiki.setAuthHandler(function(data) {
    console.info("User info", data);

    $scope.user = data;
    if ($scope.user.playlists) {
      $scope.user.playlists.forEach(function(playlist) {
        if (!playlist.items) return;
        playlist.items.sort(function(a, b) {
          return a.position - b.position;
        });
      });
    }
    $scope.views.playlistView = ($scope.getActivePlaylist() || {}).id;

    raiki.action('community.join', { slug: location.pathname.slice(1) }, function(data, err) {
      if (err) return;
      $scope.community = data;
      $('.playback .frame').attr('src', $scope.getEmbedURL());
      data.users.forEach(function(user) {
        $scope.users[user.id] = user;
      });
      $scope.$apply();
      $('#loading').fadeOut(300);
    });
  });

  raiki.on('waitlist.update', function(data) {
    $scope.community.waitlist = data;
    $scope.$apply();
  });

  raiki.on('advance', function(data) {
    console.info(data);
    $('.media .user .waitlist .text').text('Join Wait List');
    $('.media .user .waitlist.active .text').text('Leave Wait List');
    $scope.community.media = data;
    breakable: if ($scope.community.media) {
      var m = $scope.community.media;
      if (m.userID !== $scope.user.id) break breakable;
      $('.media .user .waitlist .text').text('Quit Djing');
      var playlist = $scope.getActivePlaylist();
      playlist.items.push(playlist.items.shift());
    }
    $('.playback .frame').attr('src', $scope.getEmbedURL() || '/404');
    $scope.$apply();
  });

  raiki.on('chat.receive', function(data) {
    data.now = moment();
    console.info(data);
    $scope.community.chat.push(data);
    $scope.$apply();
    $('.chat .wrapper').scrollTop($('.chat .messages').height());
  });

  raiki.on('user.join', function(data) {
    $scope.community.users.push(data);
    $scope.users[data.id] = data;
    $scope.$apply();
  });

  raiki.on('user.leave', function(data) {
    var users = $scope.community.users;
    var index = -1;
    for (var i = 0; i < users.length; i++) {
      if (users[i].id === data) {
        index = i;
        break;
      }
    }
    if (index === -1) return;
    $scope.community.users.splice(index, 1);
    $scope.$apply();
  });

  raiki.on('community.vote.woot', function(data) {
    var mehs = $scope.community.media.votes.mehs;
    var woots = $scope.community.media.votes.woots;
    for (var i = 0; i < mehs.length; i++) {
      if (mehs[i] === data) {
        mehs = mehs.splice(i, 1);
        break;
      }
    }

    for (var i = 0; i < woots.length; i++) {
      if (woots[i] === data) return;
    }

    woots.push(data);
  });

  raiki.on('community.vote.save', function(data) {
    var saves = $scope.community.media.votes.saves;
    for (var i = 0; i < saves.length; i++) {
      if (saves[i] === data) return;
    }
    saves.push(data);
  });

  raiki.on('community.vote.meh', function(data) {
    var woots = $scope.community.media.votes.woots;
    var mehs = $scope.community.media.votes.mehs;
    for (var i = 0; i < woots.length; i++) {
      if (woots[i] === data) {
        woots = woots.splice(i, 1);
        break;
      }
    }

    for (var i = 0; i < mehs.length; i++) {
      if (mehs[i] === data) return;
    }

    mehs.push(data);
  });

  raiki.connect('wss://phynix.io/_/socket', localStorage.getItem('token'));

  $('.playlist-modal .search-form').submit(function() {
    setTimeout(function() {
      var val = $(this).find('input').val();
      $scope.views.searching = true;
      $scope.views.playlistView = -1;
      $scope.$apply();
      $scope.searchYoutube(val);
    }.bind(this), 1);
  });
});

$('#chat-input').submit(function(e) {
  var input = $(this).find('.chat-input');
  var message = input.val();
  if (message === '') return;
  input.val('');
  var split = message.split('');
  var emote = false;
  if (['/em', '/me'].indexOf(split.slice(0, 3).join('')) > -1) {
    message = split.splice(3, split.length).join('');
    emote = true;
  }
  raiki.action('chat.send', { message: message, emote: emote }, function(data, err) {});
});





function importPlaylist(arr) {
  for (var i = 0; i < arr.length; i++) {
    createPlaylist(arr[i]);
  }

  function createPlaylist(p) {
    raiki.action('playlist.create', { name: p.name }, function(data, err) {
      for (var i = 0; i < p.media.length; i++) {
        insertItem(data.name, data.id, p.media[i].cid, p.media[i].format, p.media[i].title, p.media[i].artist);
      }
    });
  }

  function insertItem(name, id, cid, type, title, artist) {
    raiki.action('playlist.insert', { id: id, contentID: cid, type: type, title: title, artist: artist }, function(data, err) {
      if (err) return console.error('[%s (%s)] Failed to insert item %s - %s (%s - %s)', name, id, artist, title, cid, type);
      console.info('[%s (%s)] Successfully imported %s - %s (%s - %s [%s])', name, id, data.artist, data.title, data.contentID, data.type, data.id);
    });
  }
}