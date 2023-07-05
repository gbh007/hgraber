class Message {
  constructor(message) {
    this.close = this.close.bind(this);

    this.node = document.createElement("div");
    this.node.className = "message";
    this.node.setAttribute("level", message.level);
    this.node.innerText = message.text;
    this.node.onclick = this.close;

    if (message.autoclose) {
      setTimeout(this.close, message.autoclose * 1000);
    }
  }
  close() {
    this.node.remove();
  }
}

class Messages {
  constructor(node) {
    this.coreNode = node;

    this.handleMessage = this.handleMessage.bind(this);

    // window.dispatchEvent(new CustomEvent("new-message", { detail: {text:"text sample", level:"sample-level", autoclose: 10} }));
    window.addEventListener("new-message", this.handleMessage);
  }
  handleMessage(messageInfo) {
    console.log(messageInfo.detail);
    let msg = new Message(messageInfo.detail);
    this.coreNode.appendChild(msg.node);
  }
}

let msgs = null;

window.addEventListener("load", function () {
  let node = document.createElement("div");
  node.id = "messages";
  document.body.insertBefore(node, document.body.firstChild);
  msgs = new Messages(node);
});
