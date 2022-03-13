class Screen {
  constructor(title, xPos = null, yPos = null) {
    this.BindClose = this.BindClose.bind(this);
    this.Close = this.Close.bind(this);
    this.Node = this.Node.bind(this);

    this.node = document.createElement("div");
    this.node.className = "screen";
    if (xPos != null && yPos != null) {
      this.node.style.top = yPos + "px";
      this.node.style.left = xPos + "px";
    }

    let headNode = document.createElement("div");
    headNode.className = "screen-head";

    let titleNode = document.createElement("div");
    titleNode.className = "screen-head-title";
    titleNode.innerText = title;
    titleNode.title = title;
    titleNode.onmouseup = () => {
      this.moveOn = false;
    };
    titleNode.onmouseleave = () => {
      this.moveOn = false;
    };
    titleNode.onmousedown = () => {
      this.moveOn = true;
    };
    titleNode.onmousemove = (event) => {
      if (!this.moveOn) {
        return;
      }
      let pos = titleNode.getBoundingClientRect();
      this.node.style.top = event.pageY - (pos.bottom - pos.top) / 2 + "px";
      this.node.style.left = event.pageX - (pos.right - pos.left) / 2 + "px";
    };
    headNode.appendChild(titleNode);

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
