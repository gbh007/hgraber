class Rendering {
  calcAvg(pages) {
    if (!pages || !pages.length) {
      return 0;
    }
    return (
      Math.round(
        (pages.filter((p) => p.success).length * 10000) / pages.length
      ) / 100
    );
  }

  generateHTMLTTagArea(tags, areaName) {
    let node = document.createElement("span");
    node.style = `grid-area: ${areaName};`;
    if (tags) {
      tags.map((tagname) => {
        let tag = document.createElement("span");
        tag.className = "tag";
        tag.innerText = tagname;
        node.appendChild(tag);
      });
    }
    return node;
  }

  generateHTMLTItleDetailsFromTitleInfo(titleData, withName = true) {
    let title = document.createElement("div");
    title.className = "title-details";
    title.setAttribute("t", titleData.info.parsed.name ? "" : "bred");

    if (!titleData.pages || !titleData.pages.length) {
      let tImg = document.createElement("span");
      tImg.style = "grid-area: img;";
      title.appendChild(tImg);
    } else {
      let tImg = document.createElement("img");
      tImg.style = "max-width: 100%; max-height: 100%; grid-area: img;";
      tImg.src = `/file/${titleData.id}/1.${titleData.pages[0].ext}`;
      title.appendChild(tImg);
    }

    var node;

    if (withName) {
      node = document.createElement("span");
      node.style = "grid-area: name;";
      node.setAttribute("t", titleData.info.parsed.name ? "" : "red");
      node.innerText = titleData.info.name;
      title.appendChild(node);
    }

    node = document.createElement("span");
    node.style = "grid-area: id;";
    node.innerText = `#${titleData.id}`;
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: pgc;";
    node.setAttribute("t", titleData.info.parsed.page ? "" : "red");
    node.innerText = `Страниц: ${titleData.pages.length}`;
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: pgp;";
    node.setAttribute("t", this.calcAvg(titleData.pages) != 100.0 ? "red" : "");
    node.innerText = `Загружено: ${this.calcAvg(titleData.pages)}%`;
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: dt;";
    node.innerText = new Date(titleData.created).toLocaleString();
    title.appendChild(node);

    title.appendChild(this.generateHTMLTTagArea(titleData.info.tags, "tag"));
    title.appendChild(
      this.generateHTMLTTagArea(titleData.info.authors, "authors")
    );
    title.appendChild(
      this.generateHTMLTTagArea(titleData.info.characters, "char")
    );
    title.appendChild(
      this.generateHTMLTTagArea(titleData.info.languages, "lang")
    );
    title.appendChild(
      this.generateHTMLTTagArea(titleData.info.categories, "cat")
    );
    title.appendChild(
      this.generateHTMLTTagArea(titleData.info.parodies, "par")
    );
    title.appendChild(this.generateHTMLTTagArea(titleData.info.groups, "gr"));

    node = document.createElement("a");
    node.className = "load";
    node.innerText = "Скачать";
    title.appendChild(node);

    node = document.createElement("a");
    node.href = `/read.html?title=${titleData.id}`;
    node.className = "read";
    node.innerText = "Читать";
    title.appendChild(node);

    return title;
  }
  generateHTMLFromTitleInfo(titleData) {
    let title = document.createElement("div");
    // title.href = `/read.html?title=${info.id}`;
    title.className = "title";
    title.setAttribute("t", titleData.info.parsed.name ? "" : "bred");

    if (!titleData.pages || !titleData.pages.length) {
      let tImg = document.createElement("span");
      tImg.style = "grid-area: img;";
      title.appendChild(tImg);
    } else {
      let tImg = document.createElement("img");
      tImg.style = "max-width: 100%; max-height: 100%; grid-area: img;";
      tImg.src = `/file/${titleData.id}/1.${titleData.pages[0].ext}`;
      title.appendChild(tImg);
    }

    let node = document.createElement("span");
    node.style = "grid-area: name;";
    node.setAttribute("t", titleData.info.parsed.name ? "" : "red");
    node.innerText = titleData.info.name;
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: id;";
    node.innerText = `#${titleData.id}`;
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: pgc;";
    node.setAttribute("t", titleData.info.parsed.page ? "" : "red");
    node.innerText = `Страниц: ${titleData.pages.length}`;
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: pgp;";
    node.setAttribute("t", this.calcAvg(titleData.pages) != 100.0 ? "red" : "");
    node.innerText = `Загружено: ${this.calcAvg(titleData.pages)}%`;
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: dt;";
    node.innerText = new Date(titleData.created).toLocaleString();
    title.appendChild(node);

    node = document.createElement("span");
    node.style = "grid-area: tag;";
    (titleData.info.tags || []).map((tagname, ind) => {
      if (ind >= 8) return;
      let tag = document.createElement("span");
      tag.className = "tag";
      tag.innerText = tagname;
      node.appendChild(tag);
    });
    if (titleData.info.tags && titleData.info.tags.length > 7) {
      let more = document.createElement("b");
      more.innerText = "и больше!";
      node.appendChild(more);
    }
    title.appendChild(node);

    title.onclick = (e) => {
      let s = new Screen(titleData.info.name, e.pageX, e.pageY);
      s.Node().appendChild(
        this.generateHTMLTItleDetailsFromTitleInfo(titleData, false)
      );
      SCREENS.AddScreen(s);
    };

    return title;
  }
}
