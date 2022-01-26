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
  generateHTMLFromTitleInfo(info) {
    return `<a href="/read?title=${info.id}" class="title" t="${
      info.loaded ? "" : "bred"
    }">
      	${
          info.ext == ""
            ? '<span style="grid-area: img;"></span>'
            : '<img src="/file/' +
              info.id +
              "/1." +
              info.ext +
              '" style="max-width: 100%; max-height: 100%; grid-area: img;">'
        }
      	<span style="grid-area: name;" t="${info.loaded ? "" : "red"}">${
      info.name
    }</span>
      	<span style="grid-area: id;">#${info.id}</span>
      	<span style="grid-area: pgc;" t="${
          info.parsed_page ? "" : "red"
        }">Страниц: ${info.page_count}</span>
      	<span style="grid-area: pgp;" t="${
          info.avg != 100.0 ? "red" : ""
        }">Загружено: ${info.avg}%</span>
      	<span style="grid-area: dt;">${new Date(
          info.created
        ).toLocaleString()}</span>
      	<span style="grid-area: tag;">
      	${info.tags
          .map((tag, ind) =>
            ind < 8 ? '<span class="tag">' + tag + "</span>" : ""
          )
          .join("\n")}
             ${info.tags.length > 7 ? "<b>и больше!</b>" : ""}
      	</span>
      </a>`;
  }
}
