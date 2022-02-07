import {Article} from "./type";

type GetArticles = () => Promise<Article[]>
export const getArticles: GetArticles = async () => {
  const isString = (value: unknown): value is string => {
    return typeof value === "string"
  }
  const list = [
    ...document.querySelectorAll('.ArticleCard_container__3qUYt'),
  ];
  return Promise.all(list.map(async (element) => {
    const title = element.querySelector('.ArticleCard_title__UnBHE')?.textContent
    const timeElement = element.querySelector('time')
    const time = timeElement? timeElement.getAttribute('datetime'): null
    const urlElement = element.querySelector('.ArticleCard_mainLink__X2TOE')
    const url = urlElement? urlElement.getAttribute('href'): null
    const result = {
      title: isString(title)? title: null,
      tags: [...element.querySelectorAll('.ArticleCard_topicLink__NfdwJ')].map(
        e => e.textContent
      ).filter<string>(isString),
      time,
      url,
    };
    console.log(`selected article info: ${JSON.stringify(result)}`);
    return result
  }));
}
