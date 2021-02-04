var embed = new (function() {
  var self = this;
  var frame = $('<iframe id="embed-frame" src="https://phynix-embed.appspot.com"></iframe>');
  var buffer = [];
  var loaded = false;

  frame.on('load', function() {
    loaded = true;
    for (var i = 0; i < buffer.length; i++) {
      self.post(buffer[i]);
    }
  });

  var listeners = {
  'onVolumeChange': function() {},
  'onTimeUpdate': function() {},
  'onEnded': function() {}
  }

  this.getFrame = function() {
    return frame;
  }

  this.post = function(data) {
    if (!loaded) return buffer.push(data);
    frame.get(0).contentWindow.postMessage(JSON.stringify(data), '*');
  }

  this.destroy = function() {
    this.post({
      a: 'destroy',
      d: true
    });
  }

  this.load = function(id, type, start) {
    this.post({
      a: 'load',
      d: {
        id: id,
        type: type,
        start: start || 0
      }
    });
  }

  this.setVolume = function(volume) {
    if (volume < 0 || volume > 100) return;
    this.post({
      a: 'setVolume',
      d: volume
    });
  }

  this.onTimeUpdate = function(callback) {
    if (typeof callback !== 'function') return;
    listeners.onTimeUpdate = callback;
  }

  window.addEventListener('message', function(e) {
    var origin = e.origin || e.originalEvent.origin;
    if (origin !== 'https://phynix-embed.appspot.com') return;

    var data;
    try {
      data = JSON.parse(e.data);
    } catch (e) { return console.error(e); }

    switch (data.e) {
      case 'timeUpdate':
        listeners.onTimeUpdate(data.d);
    }
  });
});