class Api {
  async getMainInfo() {
    try {
      let response = await fetch("/info", { method: "GET" });
      return await response.json();
    } catch (err) {
      console.log(err);
      return {};
    }
  }
  async newTitle(url) {
    try {
      let response = await fetch("/new", {
        method: "POST",
        body: JSON.stringify({ url: url }),
      });
      return await response.json();
    } catch (err) {
      console.log(err);
      return {};
    }
  }
  async getTitleList(count, offset) {
    try {
      let response = await fetch("/title/list", {
        method: "POST",
        body: JSON.stringify({ count: count, offset: offset }),
      });
      return await response.json();
    } catch (err) {
      console.log(err);
      return {};
    }
  }
}
