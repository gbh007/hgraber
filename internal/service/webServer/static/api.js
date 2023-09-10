class Api {
  alert(text) {
    window.dispatchEvent(
      new CustomEvent("new-message", {
        detail: { text: text, level: "error", autoclose: 20 },
      })
    );
  }
  alertField(text) {
    window.dispatchEvent(
      new CustomEvent("new-message-field", {
        detail: { text: text, level: "error", autoclose: 20 },
      })
    );
    this.alert(text);
  }

  async getAppInfo() {
    try {
      let response = await fetch("/app/info", { method: "GET" });
      return await response.json();
    } catch (err) {
      this.alert(err);
      return {};
    }
  }

  async getMainInfo() {
    try {
      let response = await fetch("/info", { method: "GET" });
      return await response.json();
    } catch (err) {
      this.alert(err);
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
          .json()
          .then((text) => this.alertField(text))
          .catch((err) => this.alertField(err));
      } else {
        return await response.json();
      }
    } catch (err) {
      this.alertField(err);
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
      this.alert(err);
    }
    return {};
  }

  async getTitleInfo(id) {
    try {
      let response = await fetch("/title/details", {
        method: "POST",
        body: JSON.stringify({ id: id }),
      });
      return await response.json();
    } catch (err) {
      this.alert(err);
    }
    return {};
  }

  async getTitlePageInfo(id, page) {
    try {
      let response = await fetch("/title/page", {
        method: "POST",
        body: JSON.stringify({ id: id, page: page }),
      });
      return await response.json();
    } catch (err) {
      this.alert(err);
    }
    return {};
  }

  async saveToZIP(from, to) {
    try {
      let response = await fetch("/to-zip", {
        method: "POST",
        body: JSON.stringify({ from: from, to: to }),
      });
      return await response.json();
    } catch (err) {
      this.alert(err);
    }
    return {};
  }

  async updateTitleRate(id, rate) {
    return new Promise((resolve, reject) => {
      fetch("/title/rate", {
        method: "POST",
        body: JSON.stringify({ id: id, rate: rate }),
      })
        .then((response) => {
          if (!response.ok) {
            this.alert(response.statusText);
            reject();
          }
          resolve();
        })
        .catch((err) => {
          this.alert(err);
          reject();
        });
    });
  }

  async updatePageRate(id, page, rate) {
    return new Promise((resolve, reject) => {
      fetch("/title/page/rate", {
        method: "POST",
        body: JSON.stringify({ id: id, page: page, rate: rate }),
      })
        .then((response) => {
          if (!response.ok) {
            this.alert(response.statusText);
            reject();
          }
          resolve();
        })
        .catch((err) => {
          this.alert(err);
          reject();
        });
    });
  }

  async login(token) {
    try {
      let response = await fetch("/auth/login", {
        method: "POST",
        body: JSON.stringify({ token: token }),
      });
      return await response.json();
    } catch (err) {
      this.alert(err);
    }
    return {};
  }
}

const API = new Api();

class Settings {
  updateTItleOnPageCount(count) {
    let data = JSON.parse(localStorage.getItem("settings")) || {};
    data.title_on_page = parseInt(count);
    localStorage.setItem("settings", JSON.stringify(data));
    window.dispatchEvent(
      new CustomEvent("app-settings-changed", { detail: data })
    );
  }
  getTItleOnPageCount() {
    let data = JSON.parse(localStorage.getItem("settings")) || {};
    return data.title_on_page || 12;
  }
}

const SETTINGS = new Settings();
