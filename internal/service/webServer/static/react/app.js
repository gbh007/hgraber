const RouteContext = React.createContext("");
const SetRouteContext = React.createContext(console.log);

function MyApp() {
  const [route, setRoute] = React.useState("");

  return (
    <RouteContext.Provider value={route}>
      <SetRouteContext.Provider value={setRoute}>
        <TopBar />
        <Router route={route} />
      </SetRouteContext.Provider>
    </RouteContext.Provider>
  );
}

function Router(props) {
  switch (props.route) {
    case "login":
      return <Login />;

    case "books":
      return <BooksRoute />;

    default:
      return <MainRoute />;
  }
}

class MainRoute extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      workers: [],
      counts: {
        book_total: 0,
        book_not_load: 0,
        page_total: 0,
        page_not_load: 0,
      },
    };
  }
  componentDidMount() {
    this.update();
    this.updaterID = setInterval(this.update, 10000);
  }
  componentWillUnmount() {
    clearInterval(this.updaterID);
  }

  update = () => {
    MakeRequest("/info").then((data) => {
      this.setState({
        ...this.state,
        workers: data.monitor.workers,
        counts: {
          book_total: data.count,
          book_not_load: data.not_load_count,
          page_total: data.page_count,
          page_not_load: data.not_load_page_count,
        },
      });
    });
  };

  render = () => {
    return (
      <>
        <div id="messages"></div>
        <DownloadBook />
        <div id="index-info">
          <ul>
            <li>
              Всего <b>{this.state.counts.book_total}</b> тайтлов
            </li>
            <li>
              Всего незагруженно <b>{this.state.counts.book_not_load}</b>{" "}
              тайтлов
            </li>
            <li>
              Всего <b>{this.state.counts.page_total}</b> страниц
            </li>
            <li>
              Всего незагруженно <b>{this.state.counts.page_not_load}</b>{" "}
              страниц
            </li>
          </ul>
        </div>
        <DownloadArchive />

        <Workers workers={this.state.workers} />
      </>
    );
  };
}

function TopBar() {
  const setRoute = React.useContext(SetRouteContext);
  return (
    <div id="top-bar">
      <a onClick={() => setRoute("")}>Главная</a>
      <a onClick={() => setRoute("login")}>Войти</a>
      <a onClick={() => setRoute("books")}>Список тайтлов</a>
      <a onClick={() => setRoute("settings")}>Настройки</a>
    </div>
  );
}

function Workers(props) {
  return (
    <div id="info-workers">
      <table>
        <thead>
          <tr>
            <td>Название</td>
            <td>В очереди</td>
            <td>В работе</td>
            <td>Раннеров</td>
          </tr>
        </thead>
        <tbody>
          {props.workers.map((w) => (
            <tr key={w.name}>
              <td>{w.name}</td>
              <td>{w.in_queue}</td>
              <td>{w.in_work}</td>
              <td>{w.runners}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function DownloadArchive() {
  const [from, setFrom] = React.useState(0);
  const [to, setTo] = React.useState(0);

  function load() {
    MakeRequest("/to-zip", { from: from, to: to }, "POST");
  }

  return (
    <div>
      Скачать архивы С
      <input
        placeholder="С"
        value={from}
        onChange={(e) => setFrom(parseInt(e.target.value))}
      />
      По
      <input
        placeholder="По"
        value={to}
        onChange={(e) => setTo(parseInt(e.target.value))}
      />
      <button onClick={load}>Загрузить</button>
    </div>
  );
}

function DownloadBook() {
  const [errText, setErrText] = React.useState("");
  const [urlText, setURLText] = React.useState("");

  function load() {
    MakeRequest("/new", { url: urlText }, "POST")
      .then(() => {
        setURLText("");
      })
      .catch((err) => {
        setErrText(err.message);
      });
  }

  return (
    <div>
      <div>
        <input
          placeholder="Загрузить новый тайтл"
          onChange={(e) => setURLText(e.target.value)}
          value={urlText}
        />
        <span>{errText}</span>
      </div>
      <button onClick={load}>Загрузить</button>
    </div>
  );
}

function Login() {
  const [token, setToken] = React.useState("");

  function login() {
    MakeRequest("/auth/login", { token: token }, "POST");
  }

  return (
    <div>
      <input
        value={token}
        placeholder="Токен"
        onChange={(e) => setToken(e.target.value)}
      />
      <button onClick={login}>Авторизоваться</button>
    </div>
  );
}

class BooksRoute extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      books: [],
    };
  }
  componentDidMount() {
    this.refresh();
  }

  refresh = () => {
    function compressBook(book) {
      const l = book.pages.length;
      let s = 0;
      let preview = null;

      book.pages.forEach((page, index) => {
        if (!page.success) return;

        s++;

        if (!preview) {
          preview = `/file/${book.id}/${index + 1}.${page.ext}`;
        }
      });

      return {
        id: book.id,
        page_count: l,
        download: (100.0 * s) / l,
        date: book.created,
        rate: book.info.rate,
        name: book.info.name,
        tags: book.info.tags,
        preview: preview,
      };
    }
    // FIXME управление пагинацией
    MakeRequest("/title/list", { count: 30, offset: 0 }, "POST").then(
      (data) => {
        this.setState({
          ...this.state,
          books: data.map((book) => compressBook(book)),
        });
      }
    );
  };

  render = () => {
    return (
      <div id="title-list">
        {this.state.books.map((book) => (
          <TitleShortInfo key={book.id} info={book} />
        ))}
      </div>
    );
  };
}

function TitleShortInfo(props) {
  return (
    <a href="#/details.html?title=1" className="title" t="">
      <img className="title-img" src={props.info.preview} />
      <span t="" className="title-name">
        {props.info.name}
      </span>
      <span className="title-id">
        #{props.info.id}
        <Rate rate={props.info.rate} bookID={props.info.id} />
      </span>
      <span t="" className="title-page-count">
        Страниц: {props.info.page_count}
      </span>
      <span t="" className="title-page-progress">
        Загружено: {props.info.download}%
      </span>
      <span className="title-date">
        {new Date(props.info.date).toLocaleString()}
      </span>
      <span className="title-tag">
        {props.info.tags.map((tag, index) => {
          if (index > 10) return null;
          return (
            <span key={index} className="tag">
              {tag}
            </span>
          );
        })}
        {props.info.tags.length > 10 ? <b>и больше!</b> : null}
      </span>
    </a>
  );
}

function Rate(props) {
  const [rate, setRate] = React.useState(props.rate);

  function applyRate(r) {
    if (props.page) {
      MakeRequest(
        "/title/page/rate",
        { id: props.bookID, page: props.page, rate: r },
        "POST"
      ).then(() => {
        setRate(r);
      });
    } else {
      MakeRequest("/title/rate", { id: props.bookID, rate: r }, "POST").then(
        () => {
          setRate(r);
        }
      );
    }
  }

  return (
    <span className="rate">
      {[1, 2, 3, 4, 5].map((r) => (
        <span
          key={r}
          className="rate-select"
          rate={rate >= r ? rate : 0}
          onClick={() => applyRate(r)}
        ></span>
      ))}
    </span>
  );
}

async function MakeRequest(url, data = null, method = "GET") {
  return new Promise((resolve, reject) => {
    let requestInit = {
      method: method,
    };

    if (data != null) {
      requestInit.body = JSON.stringify(data);
    }

    fetch(url, requestInit)
      .then((response) => {
        // No content
        if (response.status == 204) {
          resolve(null);
          return;
        }

        if (response.ok) {
          response
            .json()
            .then((responseJSONData) => resolve(responseJSONData))
            .catch((someError) => {
              reject({ message: someError.toString() });
            });
        } else {
          response
            .json()
            .then((responseJSONData) => {
              reject({ message: responseJSONData }); // Формат ошибки строка
            })
            .catch(() => {
              response
                .text()
                .then((responseTextData) => {
                  reject({ message: responseTextData });
                })
                .catch((someError) => {
                  reject({ message: someError.toString() });
                });
            });
        }
      })
      .catch((someError) => reject({ message: someError.toString() }));
  });
}
