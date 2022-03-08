class Api {
  async getAppInfo() {
    try {
      let response = await fetch("/app/info", { method: "GET" });
      return await response.json();
    } catch (err) {
      console.log(err);
      return {};
    }
  }
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
      if (!response.ok) {
        response
          .text()
          .then((text) => alert(text))
          .catch((err) => alert(err));
      } else {
        return await response.json();
      }
    } catch (err) {
      console.log(err);
    }
    return {};
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
  async getTitleInfo(id) {
    try {
      let response = await fetch("/title/details", {
        method: "POST",
        body: JSON.stringify({ id: id }),
      });
      return await response.json();
    } catch (err) {
      console.log(err);
      return {};
    }
  }
  async getTitlePageInfo(id, page) {
    try {
      let response = await fetch("/title/page", {
        method: "POST",
        body: JSON.stringify({ id: id, page: page }),
      });
      return await response.json();
    } catch (err) {
      console.log(err);
      return {};
    }
  }
  async saveToZIP(from, to) {
    try {
      let response = await fetch("/to-zip", {
        method: "POST",
        body: JSON.stringify({ from: from, to: to }),
      });
      return await response.json();
    } catch (err) {
      console.log(err);
      return {};
    }
  }
}

const API = new Api();

class Settings {
  updateTItleOnPageCount(count) {
    let data = JSON.parse(localStorage.getItem("settings")) || {};
    data.title_on_page = parseInt(count);
    window.dispatchEvent(
      new CustomEvent("app-settings-changed", { detail: data })
    );
    localStorage.setItem("settings", JSON.stringify(data));
  }
  getTItleOnPageCount() {
    let data = JSON.parse(localStorage.getItem("settings")) || {};
    return data.title_on_page || 12;
  }
}

const SETTINGS = new Settings();
