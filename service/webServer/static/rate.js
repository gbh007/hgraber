class Rate {
  constructor(titleID, rate = null, pageID = null) {
    this.render = this.render.bind(this);
    this.updateRate = this.updateRate.bind(this);

    this.node = document.createElement("span");
    this.node.className = "rate";
    this.titleID = titleID;
    this.rate = rate || 0;
    this.pageID = pageID;

    this.render();
  }
  updateRate(rate) {
    if (this.pageID) {
      API.updatePageRate(this.titleID, this.pageID, rate).then(() => {
        this.rate = rate;
        this.render();
      });
    } else {
      API.updateTitleRate(this.titleID, rate).then(() => {
        this.rate = rate;
        this.render();
      });
    }
  }
  render() {
    this.node.innerHTML = "";

    for (let index = 1; index < 6; index++) {
      let node = document.createElement("span");
      node.className = "rate-select";
      if (this.rate >= index) {
        node.setAttribute("rate", this.rate);
      }
      node.onclick = () => this.updateRate(index);
      this.node.appendChild(node);
    }
  }
}
