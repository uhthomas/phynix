$('.playback').append(embed.getFrame());
var s;
// var app = angular.module('phynix', ['html5.sortable', 'angular-inview', 'ngJScrollPane']);
var app = angular.module('phynix', ['angular-inview', 'kyou']);
app.controller('main', ['$scope', '$timeout', '$compile', function($scope, $timeout, $compile) {
  s = $scope;
  $scope.$ = jQuery;
  $scope.moment = moment;
  $scope.$timeout = $timeout;
  $scope.settings = $.extend({
    beforeVolume: 100,
    volume: 100,
    view: {
      searchType: 1,
      itemType: 1
    }
  }, JSON.parse(localStorage.getItem('settings')) || {});
  $scope.tabIndex = 0;
  $scope.views = {
    sidebarDeployed: false,
    playlistDeployed: false,
    playlistView: 0,
    searching: false,
    searchResults: []
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
  $scope.playlists = {}
  $scope.emptyPlaylist = [];
  $scope.fullscreen = false;

  $scope.kyou = {
    playlist: {
      drop: function(e, playlist, item) {
        raiki.action('playlist.insert', { id: playlist.id, contentID: item.contentID, type: item.type }, function(data, err) {
          console.info(data, err);
          if (err) return;
          var p = $scope.getPlaylist(playlist.id);
          if (p) {
            p.items = p.items || [];
            p.items.unshift(data);
          }
          $scope.$apply();
        });
      },
      over: function(e, playlist) {

      },
      leave: function(e, playlist) {

      }
    },
    item: {
      drag: function(e, el) {
        el.css('background-image', $(e.target).find('> .image').css('background-image'));
      },
      drop: function(e, item, position) {
        if (!item.playlistID) return;
        raiki.action('playlist.item.move', { id: item.playlistID, itemID: item.id, position: position }, function(data, err) {
          if (err) return;

          item.position = position;

          var playlist = $scope.getPlaylist(item.playlistID);
          for (var i = 0; i < playlist.items.length; i++) {
            if (playlist.items[i].id === item.id) {
              playlist.items.splice(i, 1);
              break;
            }
          }

          playlist.items.splice(position, 0, item);
          $scope.$apply();
        });
      }
    }
  }

  embed.onTimeUpdate(function(data) {
    if (!$scope.community.media || !data) return;
    if (data.elapsed)
      $scope.community.media.elapsed = data.elapsed;
    if (data.buffered)
      $scope.community.media.buffered = data.buffered;
    $scope.$apply();
  });

  $scope.toggleVolume = function() {
    if (+$scope.settings.volume == 0) {
      $scope.settings.volume = +$scope.settings.beforeVolume;
    } else {
      $scope.settings.beforeVolume = +$scope.settings.volume || 100;
      $scope.settings.volume = 0;
    }
  }

  $scope.renderItem = function(e, visible) {
    var el = $(e.inViewTarget);
    var img = el.find('.image');
    if (!visible) {
      el.addClass('hidden');
      return;
    }
    el.removeClass('hidden');
  }

  $scope.sizeChat = function() {
    var jsp = $('.chat .messages').data('jsp');
    if (!jsp) return;
    jsp.reinitialise();
    $timeout(function() {
    jsp.scrollToBottom();
    }, 10);
  }

  $scope.sizePlaylistItems = function() {
    setTimeout(function() {
    $('#dynamic-styles').remove();
    var sheet = $('<style id="dynamic-styles"></style>').appendTo('head')[0];
    var payload = '';
    var areas = {
      '.playlist-modal .right': '.playlist-modal .right .items',
      '.content .tray .history': '.content .tray .history'
    };
    for (var i in areas) {
      var vars = calculate($(i).prop('scrollWidth'));
      payload += `${areas[i]} .item.grid {
        width: ${vars.width}px!important;
        height: ${vars.height}px!important;
      }`;
    }

    function calculate(width) {
      var ratio = 16 / 9;
      var min = 320;
      var items = ~~(width / min);
      width -= 30 * items;
      return {
        //width: (width / items) / width * 100,
        width: width / items,
        height: (width / items) / ratio
      };
    }

    sheet.appendChild(document.createTextNode(payload));
    }, 1);
  }

  $(window).resize(function() {
    $scope.sizeChat();
    $scope.sizePlaylistItems();
    (function resizeRoom() {
      var max = 720;
      var container = $('.content .media');
      var height = container.height();
      container.find('> .user').css('top', Math.min(max, max + (height - 74 - max)) + 'px');
    })();
  });
  $(window).resize();
  $scope.sizePlaylistItems();
  $scope.$watch('[views.playlist, views.playlistView]', function() {
    $('.playlist-modal .playlist .items .jspPane').css('top', '0px');
    $scope.sizePlaylistItems();
  });
  $scope.$watch('views.searchResults', function() {
    $('.playlist-modal .search .items .jspPane').css('top', '0px');
    $scope.sizePlaylistItems();
  });
  $scope.$watch('community.chat', function() {
    var chat = $scope.community.chat;
    var pel = $('.tray .chat .messages .jspPane');
    pel.find('.message').remove();

    var lastHead = null;
    for (var i = 0; i < chat.length; i++) {
      var c = chat[i];
      if (!c) continue;
      if (c.log) {
        $('<div class="message log"></div>').addClass(c.class).text(c.text).appendTo(pel.find('> div'));
        continue;
      }

      c.user = c.user || $scope.getUser(c.userID);
      c.now = c.now || moment(c.created);
      var msg = pel.find('> div').children().last();
      if (msg.length && +msg.attr('data-uid') === c.user.id && moment.duration(c.now.diff(lastHead.now)).minutes() <= 5) {
        var m = $('<div class="msg"></div>').text(c.message);
        if (c.emote)
          m.addClass('emote');
        m.appendTo(msg.find('.content'));
        continue;
      }

      lastHead = c;
      msg = $(`
        <div class="message">
          <div class="head">
            <div class="icons"></div>
            <div class="displayname"></div>
            <div class="username"></div>
            <div class="timestamp"></div>
          </div>
          <div class="content">
            <div class="msg"></div>
          </div>
        </div>
      `).attr('data-uid', c.user.id);
      msg.find('.displayname').text(c.user.displayname);
      msg.find('.username').text(c.user.username)
      msg.find('.timestamp').text((c.now || moment(c.created)).format('HH:mm'));
      var m = msg.find('.msg').text(c.message);

      if (c.emote)
        m.addClass('emote');

      if (c.user.id === $scope.user.id)
        msg.addClass('me');
      if ($scope.getStaff(c.user.id)) {
        msg.addClass('staff');
        var icon = {
          6: 'crown'
        }[$scope.getStaff(c.user.id).role];
        $('<i class="icon staff mdi"></i>').addClass('mdi-' + icon).appendTo(msg.find('.icons'));
      }

      msg.appendTo(pel.find('> div'));
    }
    $scope.sizeChat();
  }, true);
  $scope.$watch('settings', function() {
    localStorage.setItem('settings', JSON.stringify($scope.settings));
  }, true);
  $scope.$watch('settings.view.itemType', $scope.sizePlaylistItems);
  $scope.$watch('settings.volume', function() {
    embed.setVolume(+$scope.settings.volume);
  });

  $scope.search = function(query) {
    setTimeout(function() {
      $('.playlist-modal .search-form input').val(query);
      $scope.views.searching = true;
      $scope.views.playlist = -1;
      $scope.$apply();
      $scope.searchYoutube(query);
    }, 1);
  }

  $scope.searchYoutube = function(query) {
    $scope.views.playlistView = 1;
    $scope.views.searching = true;
    $scope.$apply();
    raiki.action('media.search', { query: query, type: 1 }, function(data, err) {
      if (err) return;
      $scope.views.searchResults = data;
      $scope.views.searchQuery = query;
      $scope.views.searching = false;
      $scope.$apply();
    });
  }

  $scope.getEmbedURL = function(item) {
    if (!$scope.community.media && !item) return;
    var preview = !!item;

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
        return `https://youtube.com/embed/${id}?autohide=1&showinfo=0&autoplay=true&controls=2&start=${~~(media.elapsed)}`;
      case 2:
        if (preview) {
          return `https://w.soundcloud.com/player/?visual=true&url=http%3A%2F%2Fapi.soundcloud.com%2Ftracks%2F${id}&show_artwork=true&auto_play=true`;
        }
        return `/s/vis.html?id=${item.contentID}&start=${~~(media.elapsed)}`;
      default:
        return '';
    }
  }

  $scope.sidebarListener = function() {
    if (!($scope.views.sidebarDeployed = !$scope.views.sidebarDeployed)) return;
    setTimeout(function() {
      $('body').one('click', function () {
        $scope.views.sidebarDeployed = false;
        $scope.$apply();
      });
    }, 1);
  }

  $scope.toggleFullscreen = function() {
    $scope.fullscreen = !(document.fullscreenElement || document.webkitFullscreenElement || document.mozFullScreenElement);
    var el = $('.media .playback').get(0);
    if ($scope.fullscreen) {
      if (el.requestFullscreen) {
        el.requestFullscreen();
      } else if (el.msRequestFullscreen) {
        el.msRequestFullscreen();
      } else if (el.mozRequestFullScreen) {
        el.mozRequestFullScreen();
      } else if (el.webkitRequestFullscreen) {
        el.webkitRequestFullscreen();
      }
      return;
    }
    if (document.cancelFullScreen) {  
      document.cancelFullScreen();  
    } else if (document.mozCancelFullScreen) {  
      document.mozCancelFullScreen();  
    } else if (document.webkitCancelFullScreen) {  
      document.webkitCancelFullScreen();  
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
    duration = ~~duration;
    var minutes = ~~(duration / 60);
    var seconds = duration % 60;
    return (minutes < 10 ? '0' + minutes : minutes) + ':' + (seconds < 10 ? '0' + seconds : seconds);
  }

  $scope.pluralize = function(num, str) {
    if (num !== 1) str += 's';
    return str;
  }

  $scope.getHost = function() {
    var staff = $scope.community.meta.staff;
    for (var i = 0; i < staff.length; i++) {
      if (staff[i].role === 6) return $scope.getUser(staff[i].userID);
    }
    return null;
  }

  $scope.getStaff = function(id) {
    var staff = $scope.community.meta.staff;
    for (var i = 0; i < staff.length; i++) {
      if (staff[i].userID === id) return staff[i];
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

  $scope.getPlaylistDuration = function(id, humanize) {
    var playlist = $scope.getPlaylist(id) || {};
    playlist.items = playlist.items || [];
    var total = 0;
    for (var i = 0; i < playlist.items.length; i++) {
      total += playlist.items[i].duration;
    }
    if (!humanize) return total;
    var pluralize = $scope.pluralize;
    var time = moment.duration(total * 1000);
    if (time.days())
      return `${time.days()} ${pluralize(time.hours(), 'day')} ${time.hours()} ${pluralize(time.hours(), 'hour')} ${time.minutes()} ${pluralize(time.minutes(), 'minute')}`;
    if (time.hours())
      return `${time.hours()} ${pluralize(time.hours(), 'hour')} ${time.minutes()} ${pluralize(time.minutes(), 'minute')}`;
    if (time.minutes())
      return `${time.minutes()} ${pluralize(time.minutes(), 'minute')} ${time.seconds()} ${pluralize(time.seconds(), 'second')}`;
    return `${time.seconds()} ${pluralize(time.seconds(), 'second')}`
  }

  $scope.getActivePlaylist = function() {
    if (!$scope.user.playlists) return;
    var playlists = $scope.user.playlists;
    for (var i = 0; i < playlists.length; i++) {
      if (playlists[i].active) return playlists[i];
    }
    return null;
  }

  $scope.playlistHasItem = function(id, cid) {
    var playlist = $scope.getPlaylist(id);
    if (!playlist) return;
    playlist.items = playlist.items || [];
    for (var i = 0; i < playlist.items.length; i++) {
      if (playlist.items[i].contentID === cid) return true;
    }
    return false;
  }

  $scope.iconActivatePlaylist = function(e) {
    var el = e.target;
    var parent = $(el).parent()
    var id = +$(el).parent().attr('data-id');
    var playlist = $scope.getActivePlaylist() || {};
    if (playlist.id === id || !parent.hasClass('viewing')) return;
    raiki.action('playlist.activate', { id: id }, function(data, err) {
      if (err) return;
      if ($scope.getActivePlaylist())
        $scope.getActivePlaylist().active = false;
      $scope.getPlaylist(id).active = true;
      $scope.$apply();
    });
  }

  $scope.woot = function() {
    if (!$scope.community.media || $scope.community.media.userID === $scope.user.id) return;
    var votes = $('.media .votes');
    if (votes.hasClass('loading')) return;
    votes.addClass('loading');
    raiki.action('community.vote.woot', {}, function(data, err) {
      votes.removeClass('loading');
    });
  }

  $scope.save = function(id) {
    if ($('.modal.add-to').length && !$('.modal.add-to').find(`[data-id=${id}]`).hasClass('invalid') && !$('.modal.add-to').find(`[data-id=${id}]`).hasClass('capacity'))  {
      $('.modal.add-to').fadeOut(100, function() {
        $('.modal.add-to').remove();
      }), void 0;
    }
    if (!$scope.community.media || $scope.community.media.userID === $scope.user.id) return;
    var votes = $('.media .votes');
    if (votes.hasClass('loading')) return;
    votes.addClass('loading');
    raiki.action('community.vote.save', { id: id }, function(data, err) {
      votes.removeClass('loading');
      if (err) return;
      var p = $scope.getPlaylist(id);
      if (p) {
        p.items = p.items || [];
        p.items.push(data);
      }
      $scope.$apply();
    });
  }

  $scope.meh = function() {
    if (!$scope.community.media || $scope.community.media.userID === $scope.user.id) return;
    var votes = $('.media .votes');
    if (votes.hasClass('loading')) return;
    votes.addClass('loading');
    raiki.action('community.vote.meh', {}, function(data, err) {
      votes.removeClass('loading');
    });
  }

  $scope.deletePlaylist = function(playlist) {
    var num = ('000' + Math.random().toString(36)).slice(-3).toUpperCase();
    var modal = $scope.compileModal('Delete Playlist', 'delete-playlist', `
      Are you sure you want to delete <span class="text"></span>
      <div class="confirmation">
        <div class="label">Enter ${num}</div>
        <form class="confirmation-form">
          <input class="input num" placeholder="${num}" maxlength="3">
        </form>
      </div>
    `, 'Delete', true);

    modal.find('.text').text(playlist.name);

    modal.find('.confirmation-form').submit(function(e) {
      e.preventDefault();
      if (modal.find('input').val() !== num) return;
      if (modal.attr('loading') === 'true') return;
      modal
        .attr('loading', true)
        .find('.message').removeClass('showing');

      raiki.action('playlist.delete', { id: playlist.id }, function(data, err) {
        modal.attr('loading', false);
        if (err) return modal.find('.message').text(err).addClass('showing');
        var p = $scope.user.playlists;
        for (var i = 0; i < p.length; i++) {
          if (p[i].id === playlist.id) {
            p.splice(i, 1);
            break;
          }
        }
        $scope.$apply();
        modal.find('.header .close').click();
      });
    });

    modal.find('.footer .submit').click(function() {
      modal.find('.confirmation-form').submit();
    });

    modal.find('input').keypress(function(e) {
      setTimeout(function() {
        console.info($(this).val());
        if ($(this).val() === num)
          return modal.find('.footer .submit').addClass('confirmed');
        modal.find('.footer .submit').removeClass('confirmed');
      }.bind(this), 1);
    });

    modal.appendTo('body').fadeIn(100, function() {
      modal.find('input').focus();
    });
  }

  $scope.editPlaylistItem = function(item) {
    var modal = $scope.compileModal('Edit Item', 'edit-item', `
      <div class="labels">
        <div class="label">Artist</div>
        <div class="label">Title</div>
      </div>
      <div class="inputs">
        <form class="edit-form">
          <input class="input artist" placeholder="Artist" maxlength="100">
          <input class="input title" placeholder="Title" maxlength="100">
          <input type="submit" style="display: none;" tabIndex="-1">
        </form>
      </div>
    `, 'Save', true);

    modal.find('.input.artist').val(item.artist);
    modal.find('.input.title').val(item.title);

    modal.find('.edit-form').submit(function(e) {
      e.preventDefault();
      if (modal.attr('loading') === 'true') return;
      modal
        .attr('loading', true)
        .find('.message').removeClass('showing');

      var artist = modal.find('.input.artist').val();
      var title = modal.find('.input.title').val();

      raiki.action('playlist.item.edit', { id: item.playlistID, itemID: item.id, artist: artist, title: title }, function(data, err) {
        modal.attr('loading', false);
        if (err) return modal.find('.message').text(err).addClass('showing');
        var p = $scope.user.playlists;
        l:
        for (var i = 0; i < p.length; i++) {
          if (p[i].id === item.playlistID) {
            for (var x = 0; x < p[i].items.length; x++) {
              if (p[i].items[x].id === item.id) {
                p[i].items[x].artist = data.artist;
                p[i].items[x].title = data.title;
                break l;
              }
            }
          }
        }
        $scope.$apply();
        modal.find('.header .close').click();
      });
    });

    modal.find('.footer .submit').click(function() {
      modal.find('.edit-form').submit();
    });

    modal.appendTo('body').fadeIn(100, function() {
      modal.find('input').focus();
    });
  }

  $scope.deletePlaylistItem = function(item) {
    if ($scope.getPlaylist(item.playlistID).items.length === 1) return;

    var modal = $scope.compileModal('Delete Item', 'delete-item', `
      Are you sure you want to delete
      <div class="text"></div>
      from your playlist?
    `, 'Delete', true);

    modal.find('.content .text').text(item.artist + ' - ' + item.title);

    modal.find('.footer .submit').click(function() {
      if (modal.attr('loading') === 'true') return;
      modal
        .attr('loading', true)
        .find('.message').removeClass('showing');

      raiki.action('playlist.item.delete', { id: item.playlistID, itemID: item.id }, function(data, err) {
        console.info(data, err);
        modal.attr('loading', false);
        if (err) return modal.find('.message').text(err).addClass('showing');
        var p = $scope.user.playlists;
        l:
        for (var i = 0; i < p.length; i++) {
          if (p[i].id === item.playlistID) {
            for (var x = 0; x < p[i].items.length; x++) {
              if (p[i].items[x].id === item.id) {
                p[i].items.splice(x, 1);
                break l;
              }
            }
          }
        }
        $scope.$apply();
        modal.find('.header .close').click();
      });
    });

    modal.appendTo('body').fadeIn(100, function() {
      modal.find('input').focus();
    });
  }

  $scope.itemPreview = function(item) {
    var before = $scope.settings.volume;
    var url = $scope.getEmbedURL(item);
    var modal = $scope.compileModal(item.artist + ' - ' + item.title, 'item-preview', `
      <iframe class="frame" src="${url}" frameborder="0" allowfullscreen></iframe>
    `, 'Close', false);

    modal.find('.footer .submit').click(function() {
      modal.fadeOut(100, function() {
        this.remove();
      })
    });

    modal.click(function(e) {
      if (e.target !== this) return;
      $scope.settings.volume = before;
      $scope.$apply();
    });
    modal.find('.footer .submit, .footer .cancel, .header .mdi-close').click(function() {
      $scope.settings.volume = before;
      $scope.$apply();
    });

    modal.appendTo('body').fadeIn(100);

    $scope.settings.volume = 0;
  }

  $scope.showAddToModal = function(save) {
    if ($('.modal.add-to').length)  {
      return $('.modal.add-to').fadeOut(100, function() {
        $('.modal.add-to').remove();
      }), void 0;
    }

    if (save && !$scope.community.media || $scope.community.media.userID === $scope.user.id) return;

    var txt = save ? 'save' : 'add';
    var icon = save ? 'star' : 'plus';
    var modal = $scope.compileModal(`${txt} To`, 'add-to', `
      <div class="playlist" ng-repeat="playlist in user.playlists | orderBy:'name'" ng-class="{'active': playlist.active, 'invalid': playlistHasItem(playlist.id, community.media.item.contentID), 'capacity': playlist.items.length >= 500}" ng-click="save(playlist.id)" data-id="{{ playlist.id }}">
        <i class="icon mdi mdi-check" ng-class="{'mdi-close': playlistHasItem(playlist.id, community.media.item.contentID) || playlist.items.length >= 500}"></i>
        <div class="name">{{ playlist.name }}</div>
        <div class="count">{{ playlist.items.length || 0 }}</div>
      </div>
      `, '', false);
    modal.find('.content').attr('kyou-container', '');
    modal.find('.loader').remove();
    modal.find('.footer').remove();

    $(`
      <div class="icon">
        <i class="mdi mdi-close"></i>
        <i class="mdi mdi-${icon}"></i>
      </div>
    `).click(function() {
      $(modal).fadeOut(100, function() {
        modal.remove();
      });
    }).prependTo(modal.find('.header'));

    $($compile(modal)($scope)).appendTo('.content .media .user').fadeIn(100);
  }

  $scope.showPlaylistCreateModal = function(item) {
    var modal = $scope.compileModal('Create Playlist', 'create-playlist', `
      <form class="playlist-create-form">
        <input class="name-input" maxlength="30" placeholder="Name your playlist">
        <input type="submit" style="display: none;" tabIndex="-1">
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

        $scope.user.playlists = $scope.user.playlists || [];
        $scope.user.playlists.map(function(playlist) {
          playlist.active = false;
        });
        $scope.user.playlists.push(data);
        $scope.views.playlistView = data.id;
        $scope.$apply();
        modal.find('.header .close').click();
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
            <i class="close mdi mdi-close"></i>
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
    $scope.views.playlist = ($scope.getActivePlaylist() || {}).id;

    raiki.action('community.join', { slug: location.pathname.slice(1) }, function(data, err) {
      if (err) return;
      $scope.community = data;
      data.chat = data.chat || [];
      data.users = data.users || [];
      data.waitlist = data.waitlist || [];
      data.chat.forEach(function(chat) {
        $scope.users[chat.user.id] = chat.user;
      });
      data.users.forEach(function(user) {
        $scope.users[user.id] = user;
      });

      if (data.media) {
        embed.setVolume(+$scope.settings.volume);
        embed.load(data.media.item.contentID, data.media.item.type, data.media.elapsed);
      }

      $scope.$apply();
      $('#loading').fadeOut(300);
    });
  });

  raiki.on('waitlist.update', function(data) {
    $scope.community.waitlist = data;
    $scope.$apply();
  });

  raiki.on('advance', function(data) {
    $('.media .user .waitlist .text').text('Join Wait List');
    $('.media .user .waitlist.active .text').text('Leave Wait List');
    if (data.history) {
      $scope.community.meta.history.unshift(data.history);
      if ($scope.community.meta.history.length > 50) {
        $scope.community.meta.history.splice(50, $scope.community.meta.history.length - 50);
      }
    }
    $scope.community.media = data.media;
    breakable: if ($scope.community.media) {
      var m = $scope.community.media;
      if (m.userID !== $scope.user.id) break breakable;
      $('.media .user .waitlist .text').text('Quit DJing');
      var playlist = $scope.getActivePlaylist();
      playlist.items.push(playlist.items.shift());
    }
    if (data.media) {
      embed.setVolume(+$scope.settings.volume);
      embed.load(data.media.item.contentID, data.media.item.type, 0);
    } else {
      embed.destroy();
    }
    $scope.$apply();
  });

  raiki.on('chat.receive', function(data) {
    data.now = moment();
    console.info(data);
    $scope.community.chat.push(data);
    if ($scope.community.chat.length > 255) {
      $scope.community.chat.splice(0, $scope.community.chat.length - 255);
    }
    $scope.$apply();
  });

  raiki.on('user.join', function(data) {
    $scope.community.users.push(data);
    $scope.users[data.id] = data;
    $scope.community.chat.push({ log: true, text: `${data.displayname} joined the community`, class: 'join' })
    $scope.$apply();
    $('.chat .wrapper').scrollTop($('.chat .messages').height());
  });

  raiki.on('user.leave', function(data) {
    var users = $scope.community.users;
    for (var i = 0; i < users.length; i++) {
      if (users[i].id !== data) continue;
      $scope.community.chat.push({ log: true, text: `${users[i].displayname} left the community`, class: 'leave' })
      $scope.community.users.splice(i, 1);
      break;
    }

    var media = $scope.community.media;
    if (media) {
      var woots = media.votes.woots;
      var mehs = media.votes.mehs;
      if (woots.indexOf(data) > -1) {
        woots.splice(woots.indexOf(data), 1);
      }
      if (mehs.indexOf(data) > -1) {
        mehs.splice(mehs.indexOf(data), 1);
      }
    }
    $scope.$apply();
    $('.chat .wrapper').scrollTop($('.chat .messages').height());
  });

  raiki.on('community.vote.woot', function(data) {
    var mehs = $scope.community.media.votes.mehs;
    var woots = $scope.community.media.votes.woots;
    if (mehs.indexOf(data) > -1) {
      mehs.splice(mehs.indexOf(data), 1);
    }

    if (woots.indexOf(data) === -1) {
      woots.push(data);
    }

    $scope.$apply();
  });

  raiki.on('community.vote.save', function(data) {
    var saves = $scope.community.media.votes.saves;
    if (saves.indexOf(data) === -1) {
      saves.push(data);
    }
    $scope.$apply();
  });

  raiki.on('community.vote.meh', function(data) {
    var woots = $scope.community.media.votes.woots;
    var mehs = $scope.community.media.votes.mehs;
    if (woots.indexOf(data) > -1) {
      woots.splice(woots.indexOf(data), 1);
    }

    if (mehs.indexOf(data) === -1) {
      mehs.push(data);
    }

    $scope.$apply();
  });

  raiki.connect('wss://phynix.io/_/socket', localStorage.getItem('token'));

  $('.playlist-modal .search-form').submit(function() {
    setTimeout(function() {
      var val = $(this).find('input').val();
      $scope.search(val);
    }.bind(this), 1);
  });
}]);

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