class Messages {
  constructor(node) {
    this.coreNode = node;
    this.messages = [];

    let self = this;

    // window.dispatchEvent(new CustomEvent("new-message", { detail: {text:"text sample"} }));
    window.addEventListener("new-message", function (data) {
      self.handleMessage(data);
    });
  }
  handleMessage(messageInfo) {
    console.log(messageInfo.detail);
    this.messages.push(messageInfo.detail);
  }
}

let msgs = null;

window.addEventListener("load", function () {
  msgs = new Messages(document.getElementById("top-bar"));
});
