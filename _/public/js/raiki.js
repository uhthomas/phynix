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