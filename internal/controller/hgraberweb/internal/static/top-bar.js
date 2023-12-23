window.addEventListener("load", function () {
  let node = document.createElement("div");
  node.id = "top-bar";
  node.innerHTML = `<a href="/">Главная</a>
<a href="/auth.html">Войти</a>
<a href="/list.html">Список тайтлов</a>
<a href="/settings.html">Настройки</a>`;
  document.body.insertBefore(node, document.body.firstChild);
});
