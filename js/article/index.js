'use strict';

document.addEventListener('DOMContentLoaded', () => {
  const deleteBtns = document.querySelectorAll('.articles__item-delete');
  const moreBtn = document.querySelector('.page__more');
  const articles = document.querySelector('.articles');
  const articleTmpl = document.querySelector('.articles__item-tmpl').firstElementChild;

  const csrfToken = document.getElementsByName('csrf')[0].content;
  const deleteArticle = id => {
    let statusCode;

    fetch(`/api/articles/${id}`, {
      method: 'DELETE',
      headers: { 'X-CSRF-Token': csrfToken }
    })
      .then(res => {
        statusCode = res.status;
        return res.json();
      })
      .then(data => {
        console.log(JSON.stringify(data));
        if (statusCode == 200) {
          document.querySelector(`.articles__item-${id}`).remove();
        }
      })
      .catch(err => console.error(err));
  };

  for (let elm of deleteBtns) {
    elm.addEventListener('click', event => {
      event.preventDefault();
      deleteArticle(elm.dataset.id);
    });
  }

  moreBtn.addEventListener('click', event => {
    event.preventDefault();
    const cursor = moreBtn.dataset.cursor;
    if (!cursor || cursor <= 0) {
      moreBtn.remove();
      return;
    }

    let statusCode;
    fetch(`/api/articles?cursor=${cursor}`)
      .then(res => {
        statusCode = res.status;
        return res.json();
      })
      .then(data => {
        console.log(JSON.stringify(data));
        if (statusCode == 200 && Array.isArray(data)) {
          if (data.length == 0) {
            moreBtn.remove();
            return;
          }

          const fragment = document.createDocumentFragment();
          data.forEach(article => {
            const frag = document.createDocumentFragment();
            frag.appendChild(articleTmpl.cloneNode(true));
            frag.querySelector('article').classList.add(`articles__item-${article.id}`);
            frag.querySelector('.articles__item').setAttribute('href', `/articles/${article.id}`);
            frag.querySelector('.articles__item-title').textContent = article.title;
            frag.querySelector('.articles__item-date').textContent = article.created.split('T')[0]; 

            const deleteBtnElm = frag.querySelector('.articles__item-delete');
            deleteBtnElm.dataset.id = article.id;
            deleteBtnElm.addEventListener('click', event => {
              event.preventDefault();
              deleteArticle(article.id);
            });

            fragment.appendChild(frag);
          });

          moreBtn.dataset.cursor = data[data.length - 1].id;

          articles.appendChild(fragment);
        }
      })
      .catch(err => console.error(err));
  });
  bubbly({
    colorStart: '#fff4e6',
    colorStop: '#ffe9e4',
    blur: 1,
    compose: 'source-over',
    bubbleFunc:() => `hsla(${Math.random() * 50}, 100%, 50%, .3)`
  });

});
