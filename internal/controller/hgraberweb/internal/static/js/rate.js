class Rate {
  constructor(node) {
    this.render = this.render.bind(this);
    this.updateRate = this.updateRate.bind(this);

    document.createElement("span").getAttribute("");

    this.node = node;

    this.bookID = parseInt(this.node.getAttribute("book"));
    this.rate = parseInt(this.node.getAttribute("rate") || "0");
    this.pageNumber = parseInt(this.node.getAttribute("page") || "0");

    this.render();
  }
  async updateRate(rate) {
    try {
      let response = await fetch("/api/rate", {
        method: "POST",
        body: JSON.stringify({
          id: this.bookID,
          page: this.pageNumber || 0,
          rate: rate,
        }),
      });

      if (!response.ok) {
        throw new Error(await response.text());
      }

      this.rate = rate;
      this.render();
    } catch (e) {
      console.log(e);
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

function refreshRates() {
  window.dispatchEvent(new Event("app-refresh-rates"));
}

window.addEventListener("app-refresh-rates", function () {
  document.querySelectorAll("span.rate[unprocessed]").forEach((node) => {
    new Rate(node);
    node.removeAttribute("unprocessed");
  });
});
