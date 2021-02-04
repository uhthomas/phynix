/*
  Dependencies:
  * JQuery
  * Youtube Player API 
*/

var Playback = function(el) {
  var types = {
    youtube: 1,
    soudncloud: 2
  }
  var volume = 1;
  var paused = false;
  var frame = $(el);
  var mediaType;
  var media; // Object -- Youtube Player Object or Phynix soundcloud player

  this.pause = function() {

  }

  this.play = function() {
    if (!media) return;
    switch (mediaType) {

    }
  }

  this.destroy = function() {
    if (!media) return;
    switch (mediaType) {
      case types.youtube:
        media.destroy();
        break;
      case types.soundcloud:
        this.pause();
        break;
    }
    media = void 0;
    mediaType = void 0;
  }

  // this.queue

  this.setVolume = function (value) {
    if (value < 0 || value > 1) return;
    switch (mediaType) {
      case types.youtube:
        media.setVolume(value * 100);
    }
  }

  this.setMedia = function(id, type, duration) {
    if (!!media && type !== types.youtube) this.destroy();
    switch (type) {
      case types.youtube:
        if (media instanceof YT.Player) return media.loadVideoById(id, duration || 0);

        media = new YT.Player(frame.attr('id'), {
          width: '100%',
          height: '100%',
          videoId: id,
          playerVars: {
            'autoplay': 1,
            'controls': 1
          },
          events: {
            'onReady': function(e) {
              console.info(e);
              e.target.playVideo();
            },
            'onStateChange': function(e) {
              switch (e.data) {
                case YT.PlayerState:
              }
            }
          }
        });
        break;
      case types.soundcloud:

    }
  }

  this.setQuality = function(quality) {

  }

  this.onTick = function(callback) {
    
  }
}