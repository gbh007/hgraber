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
class Rendering {
  generateHTMLTItleDetailsFromTitleInfo(info) {
    let title = document.createElement("div");
    title.className = "title-details";
    title.setAttribute("t", info.loaded ? "" : "bred");

    if (info.ext == "") {
      let tImg = document.createElement("span");
      tImg.style = "grid-area: img;";
      title.appendChild(tImg);
    } else {
      let tImg = document.createElement("img");
      tImg.style = "max-width: 100%; max-height: 100%; grid-area: img;";
      tImg.src = `/file/${info.id}/1.${info.ext}`;
      title.appendChild(tImg);
    }

    let node = document.createElement("span");
    node.style = "grid-area: name;";
    node.setAttribute("t", info.loaded ? "" : "red");
    node.innerText = info.name;
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: id;";
    node.innerText = `#${info.id}`;
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: pgc;";
    node.setAttribute("t", info.parsed_page ? "" : "red");
    node.innerText = `Страниц: ${info.page_count}`;
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: pgp;";
    node.setAttribute("t", info.avg != 100.0 ? "red" : "");
    node.innerText = `Загружено: ${info.avg}%`;
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: dt;";
    node.innerText = new Date(info.created).toLocaleString();
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: tag;";
    info.tags.map((tagname) => {
      let tag = document.createElement("span");
      tag.className = "tag";
      tag.innerText = tagname;
      node.appendChild(tag);
    });
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: authors;";
    info.authors.map((tagname) => {
      let tag = document.createElement("span");
      tag.className = "tag";
      tag.innerText = tagname;
      node.appendChild(tag);
    });
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: char;";
    info.characters.map((tagname) => {
      let tag = document.createElement("span");
      tag.className = "tag";
      tag.innerText = tagname;
      node.appendChild(tag);
    });
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: lang;";
    info.languages.map((tagname) => {
      let tag = document.createElement("span");
      tag.className = "tag";
      tag.innerText = tagname;
      node.appendChild(tag);
    });
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: cat;";
    info.categories.map((tagname) => {
      let tag = document.createElement("span");
      tag.className = "tag";
      tag.innerText = tagname;
      node.appendChild(tag);
    });
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: par;";
    info.parodies.map((tagname) => {
      let tag = document.createElement("span");
      tag.className = "tag";
      tag.innerText = tagname;
      node.appendChild(tag);
    });
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: gr;";
    info.groups.map((tagname) => {
      let tag = document.createElement("span");
      tag.className = "tag";
      tag.innerText = tagname;
      node.appendChild(tag);
    });
    title.appendChild(node);

    node = document.createElement("a");
    node.className = "load";
    node.innerText = "Скачать";
    title.appendChild(node);

    node = document.createElement("a");
    node.href = `/read?title=${info.id}`;
    node.className = "read";
    node.innerText = "Читать";
    title.appendChild(node);

    return title;
  }
  generateHTMLFromTitleInfo(info) {
    let title = document.createElement("div");
    // title.href = `/read?title=${info.id}`;
    title.className = "title";
    title.setAttribute("t", info.loaded ? "" : "bred");

    if (info.ext == "") {
      let tImg = document.createElement("span");
      tImg.style = "grid-area: img;";
      title.appendChild(tImg);
    } else {
      let tImg = document.createElement("img");
      tImg.style = "max-width: 100%; max-height: 100%; grid-area: img;";
      tImg.src = `/file/${info.id}/1.${info.ext}`;
      title.appendChild(tImg);
    }

    let node = document.createElement("span");
    node.style = "grid-area: name;";
    node.setAttribute("t", info.loaded ? "" : "red");
    node.innerText = info.name;
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: id;";
    node.innerText = `#${info.id}`;
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: pgc;";
    node.setAttribute("t", info.parsed_page ? "" : "red");
    node.innerText = `Страниц: ${info.page_count}`;
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: pgp;";
    node.setAttribute("t", info.avg != 100.0 ? "red" : "");
    node.innerText = `Загружено: ${info.avg}%`;
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: dt;";
    node.innerText = new Date(info.created).toLocaleString();
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: tag;";
    info.tags.map((tagname, ind) => {
      if (ind >= 8) return;
      let tag = document.createElement("span");
      tag.className = "tag";
      tag.innerText = tagname;
      node.appendChild(tag);
    });
    if (info.tags.length > 7) {
      let more = document.createElement("b");
      more.innerText = "и больше!";
      node.appendChild(more);
    }
    title.appendChild(node);

    title.appendChild(this.generateHTMLTItleDetailsFromTitleInfo(info));

    return title;
  }
}
