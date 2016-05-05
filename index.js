$(function () {
  var ws = new WebSocket("ws://localhost:8080/connection");
  var $btn = $('.js-btn-submit');
  var $text = $('.js-text-message');
  var $container = $('.js-messages');

  $btn.on('click', function () {
    var requestJson = {
      text:  $text.val()
    };
    ws.send(JSON.stringify(requestJson));
  });

  ws.onmessage = function (e) {
    var jsonData = JSON.parse(e.data);
    $container.append("<p>" + jsonData.Text + "</p>");
  }
});
