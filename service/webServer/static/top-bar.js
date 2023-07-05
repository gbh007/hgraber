class TopBar {
  constructor(node) {
    this.coreNode = node;

    this.goMainPage = this.goMainPage.bind(this);
    this.showSettings = this.showSettings.bind(this);
    this.hideSettings = this.hideSettings.bind(this);
    this.goToTitles = this.goToTitles.bind(this);

    // Главная
    this.mainLink = document.createElement("span");
    this.mainLink.innerText = "Главная";
    this.mainLink.onclick = this.goMainPage;
    this.coreNode.appendChild(this.mainLink);

    // Список тайтлов
    this.titleListLink = document.createElement("span");
    this.titleListLink.innerText = "Список тайтлов";
    this.titleListLink.onclick = this.goToTitles;
    this.coreNode.appendChild(this.titleListLink);

    // Настройки
    this.settingsLink = document.createElement("span");
    this.settingsLink.innerText = "Настройки";
    this.settingsLink.onclick = this.showSettings;
    this.coreNode.appendChild(this.settingsLink);

    this.settingsNode = document.createElement("div");
    this.settingsNode.className = "top-bar-settings";
    this.coreNode.appendChild(this.settingsNode);
  }
  goMainPage() {
    console.log(document.location.href);
    document.location.href = "/";
  }
  goToTitles() {
    console.log(document.location.href);
    document.location.href = "/list.html";
  }
  showSettings() {
    this.renderSettings().then(() => {
      this.settingsNode.setAttribute("show", "1");
    });
  }
  hideSettings() {
    this.settingsNode.setAttribute("show", "0");
  }
  async renderSettings() {
    this.settingsNode.querySelectorAll("*").forEach((n) => n.remove());

    let info = await API.getAppInfo();

    let table = document.createElement("table");
    let tbody = document.createElement("tbody");
    let tr, td;

    tr = document.createElement("tr");
    td = document.createElement("td");
    td.innerText = "Данные сервера";
    td.className = "head";
    td.setAttribute("colspan", "2");
    tr.appendChild(td);
    tbody.appendChild(tr);

    tr = document.createElement("tr");
    td = document.createElement("td");
    td.innerText = "Версия";
    tr.appendChild(td);
    td = document.createElement("td");
    td.innerText = info.version;
    tr.appendChild(td);
    tbody.appendChild(tr);

    tr = document.createElement("tr");
    td = document.createElement("td");
    td.innerText = "Коммит";
    tr.appendChild(td);
    td = document.createElement("td");
    td.innerText = info.commit;
    tr.appendChild(td);
    tbody.appendChild(tr);

    tr = document.createElement("tr");
    td = document.createElement("td");
    td.innerText = "Время сборки";
    tr.appendChild(td);
    td = document.createElement("td");
    td.innerText = info.build_at;
    tr.appendChild(td);
    tbody.appendChild(tr);

    tr = document.createElement("tr");
    td = document.createElement("td");
    td.innerText = "Данные приложения";
    td.className = "head";
    td.setAttribute("colspan", "2");
    tr.appendChild(td);
    tbody.appendChild(tr);

    tr = document.createElement("tr");
    td = document.createElement("td");
    td.innerText = "Количество на странице";
    tr.appendChild(td);
    td = document.createElement("td");
    let input = document.createElement("input");
    input.setAttribute("type", "number");
    input.value = SETTINGS.getTItleOnPageCount();
    input.oninput = (v) => SETTINGS.updateTItleOnPageCount(v.target.value);
    td.appendChild(input);
    tr.appendChild(td);
    tbody.appendChild(tr);

    tr = document.createElement("tr");
    td = document.createElement("td");
    td.setAttribute("colspan", "2");
    let button = document.createElement("button");
    button.innerText = "Закрыть";
    button.onclick = this.hideSettings;
    td.appendChild(button);
    tr.appendChild(td);
    tbody.appendChild(tr);

    table.appendChild(tbody);
    this.settingsNode.appendChild(table);
  }
}

let topBar = null;

window.addEventListener("load", function () {
  let node = document.createElement("div");
  node.id = "top-bar";
  document.body.insertBefore(node, document.body.firstChild);
  topBar = new TopBar(node);
});
