(function($) {
  var module = angular.module('kyou', []);

  module.directive('kyouDropzone', ['$timeout', function($timeout) {
    function link(scope, el, attrs) {
      var kyou = {
        last: null
      };

      kyou.drop = function(e) {
        e.originalEvent.preventDefault();
        $(el).removeClass(scope.overClass);
        scope.onDrop()(e, scope.data, JSON.parse(e.originalEvent.dataTransfer.getData('Text')).data);
      }

      kyou.dragenter = function(e) {
        e.originalEvent.preventDefault();
        kyou.last = e.originalEvent.target;
        if (scope.overClass) $(el).addClass(scope.overClass);
        scope.onOver()(e, scope.data);
        e.preventDefault()
      }

      kyou.dragleave = function(e) {
        e.originalEvent.preventDefault();
        if (kyou.last !== e.originalEvent.target) return;
        $(el).removeClass(scope.overClass);
        scope.onLeave()(e, scope.data);
      }

      kyou.dragover = function(e) {
        e.preventDefault();
        scope.onOver()(e, scope.data);
      }

      $(el)
        .on('drop', kyou.drop)
        .on('dragenter', kyou.dragenter)
        .on('dragleave', kyou.dragleave)
        .on('dragover', kyou.dragover);
    }

    return {
      restrict: 'A',
      scope: {
        onDrop: '&kyouDropzoneOnDrop',
        onOver: '&kyouDropzoneOnOver',
        onLeave: '&kyouDropzoneOnLeave',
        overClass: '@kyouDropzoneOverClass',
        data: '=kyouDropzoneData'
      },
      link: link
    }
  }]);

  $(document).on('dragover', function(e) {
    if (!kyouItem) return;
    kyouItem.drag(e);
  });

  var kyouitem = null;
  module.directive('kyouItem', function() {
    var canvas = document.createElement('canvas');
    canvas.width = 1;
    canvas.height = 1;
    var image = new Image();
    image.src = canvas.toDataURL();

    function link(scope, el, attrs) {
      $(el).attr('draggable', true);

      var kyou = {
        dragging: false,
        dragItem: null,
        dragItemX: 0,
        dragItemY: 0
      };

      kyou.updateDragItem = function() {
        if (!kyou.dragging) return;
        window.requestAnimationFrame(kyou.updateDragItem);
        // kyou.dragItem.css({
        //   left: kyou.dragItemX,
        //   top: kyou.dragItemY
        // });
        // kyou.dragItem.css('transform', `translate3d(${kyou.dragItemX}px, ${kyou.dragItemY}px, 99px)`);
        kyou.dragItem.get(0).style.transform = `translate3d(${kyou.dragItemX}px, ${kyou.dragItemY}px, 99px)`;
      }

      kyou.dragstart = function(e) {
        e = e.originalEvent;
        e.dataTransfer.effectAllowed = 'move';
        e.dataTransfer.setDragImage(image, 1, 1);
        e.dataTransfer.setData('Text', JSON.stringify({position: $(e.target).index()-1, data: scope.data}));

        if (scope.dragClass)
          $(el).addClass(scope.dragClass);

        kyouItem = kyou;
        kyou.dragItem = $('<div class="drag-item"></div>').appendTo('body');
        kyou.dragging = true;
        kyou.updateDragItem();

        scope.onDrag()(e, kyou.dragItem);
      }

      kyou.drag = function(e) {
        e = e.originalEvent;
        var x = e.clientX || e.pageX;
        var y = e.clientY || e.pageY;
        if (!x || !y) return;
        kyou.dragItemX = x;
        kyou.dragItemY = y;
      }

      kyou.dragend = function(e) {
        e = e.originalEvent;
        $('.drag-item').remove();
        $(e).removeClass(scope.overClass);
        kyouItem = null;
        kyou.dragging = false;
        kyou.dragItem = null;
      }

      kyou.drop = function(e) {
        e.originalEvent.preventDefault();
        $(el).removeClass(scope.overClass);

        var data = JSON.parse(e.originalEvent.dataTransfer.getData('Text'));

        var top = e.offsetY <= $(el).height()/2;
        var left = e.offsetX <= $(el).width()/2;
        var position = $(el).index()-2;

        if (scope.sidePosition && !left) {
          position++;
        } else if (!scope.sidePosition && !top) {
          position++;
        }

        if (position < data.position) {
          position++;
        }

        scope.onDrop() && scope.onDrop()(e, data.data, Math.max(0, position));
      }

      kyou.dragenter = function(e) {
        e.originalEvent.preventDefault();
        kyou.last = e.originalEvent.target;
        if (scope.overClass) $(el).addClass(scope.overClass);
        scope.onOver() && scope.onOver()(e, scope.data);
      }

      kyou.dragleave = function(e) {
        e.originalEvent.preventDefault();
        if (kyou.last !== e.originalEvent.target) return;
        $(el).removeClass(scope.overClass);
        scope.onLeave() && scope.onLeave()(e, scope.data);
      }

      kyou.dragover = function(e) {
        e.preventDefault();
        scope.onOver() && scope.onOver()(e, scope.data);
      }

      $(el)
        .on('drop', kyou.drop)
        .on('dragstart', kyou.dragstart)
        .on('drag', kyou.drag)
        .on('dragend', kyou.dragend)
        .on('dragenter', kyou.dragenter)
        .on('dragleave', kyou.dragleave)
        .on('dragover', kyou.dragover);
    }

    return {
      restrict: 'A',
      scope: {
        overClass: '@kyouItemOverClass',
        dragClass: '@kyouItemDragClass',
        onDrag: '&kyouItemOnDrag',
        onDrop: '&kyouItemOnDrop',
        onOver: '&kyouItemOnOver',
        onLeave: '&kyouItemOnLeave',
        sidePosition: '=kyouItemSidePosition',
        data: '=kyouItem'
      },
      link: link
    }
  });

  module.directive('kyouContainer', function() {
    function link(scope, el, attrs) {
      $(el).jScrollPane({ autoReinitialise: true });
      var kyou = {};

      kyou.drop = function(e) {
        e.preventDefault();
      }

      kyou.dragenter = function(e) {
        e.preventDefault();
      }

      kyou.dragover = function(e) {
        e.preventDefault();
      }

      kyou.dragleave = function(e) {
        e.preventDefault();
      }

      $(el)
        .on('drop', kyou.drop)
        .on('dragenter', kyou.dragenter)
        .on('dragover', kyou.dragover)
        .on('dragleave', kyou.dragleave);
    }

    function compile(el, attrs, transclude) {
      console.info(arguments);
      return link;
    }

    return {
      restrict: 'A',
      scope: {
        onDrop: '&kyouContainerOnDrop',
        onOver: '&kyouContainerOnOver',
        onLeave: '&kyouContainerOnLeave',
        overClass: '@kyouContainerOverClass',
      },
      transclude: true,
      template: function(el, attrs) {
        return '<div class="scroll-pane"><div ng-transclude></div></div>';
      },
      replace: true,
      link: link
    }
  });
})(jQuery);