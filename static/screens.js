class Screen {
  constructor(title, xPos = null, yPos = null) {
    this.BindClose = this.BindClose.bind(this);
    this.Close = this.Close.bind(this);
    this.Node = this.Node.bind(this);
    this.move = this.move.bind(this);

    this.node = document.createElement("div");
    this.node.className = "screen";
    if (xPos != null && yPos != null) {
      this.node.style.top = yPos + "px";
      this.node.style.left = xPos + "px";
    }

    let headNode = document.createElement("div");
    headNode.className = "screen-head";

    this.titleNode = document.createElement("div");
    this.titleNode.className = "screen-head-title";
    this.titleNode.innerText = title;
    this.titleNode.title = title;
    this.titleNode.onmouseup = () => {
      window.removeEventListener("mousemove", this.move);
    };

    this.titleNode.onmousedown = () => {
      window.addEventListener("mousemove", this.move);
    };

    headNode.appendChild(this.titleNode);

    let actionsNode = document.createElement("div");
    actionsNode.className = "screen-head-actions";

    let actionCloseNode = document.createElement("span");
    actionCloseNode.className = "screen-head-actions-close";
    actionCloseNode.onclick = this.Close;
    actionsNode.appendChild(actionCloseNode);

    headNode.appendChild(actionsNode);

    this.node.appendChild(headNode);

    this.contentNode = document.createElement("div");
    this.contentNode.className = "screen-content";
    this.node.appendChild(this.contentNode);
  }
  move(event) {
    let pos = this.titleNode.getBoundingClientRect();
    this.node.style.top = event.pageY - (pos.bottom - pos.top) / 2 + "px";
    this.node.style.left = event.pageX - (pos.right - pos.left) / 2 + "px";
  }
  Close() {
    if (this.close) {
      this.close();
    }
    this.node.remove();
  }
  BindClose(close = null) {
    this.close = close;
  }
  Node() {
    return this.contentNode;
  }
}

class ScreenController {
  constructor(node) {
    this.AddScreen = this.AddScreen.bind(this);

    this.coreNode = node;
  }
  AddScreen(screen) {
    this.coreNode.appendChild(screen.node);
  }
}

let SCREENS = null;

window.addEventListener("load", function () {
  let node = document.createElement("div");
  node.id = "screens";
  document.body.insertBefore(node, document.body.firstChild);
  SCREENS = new ScreenController(node);
});
