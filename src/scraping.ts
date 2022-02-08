import puppeteer from 'puppeteer';
import {Article} from "./type";
import {getArticles} from "./getArticles";
import {omit} from "remeda";

const url = 'https://zenn.dev';

type Start = (name: string) => Promise<Article[]>
export const start: Start = async (name) => {
  console.log('start');
  const browser = await puppeteer.launch({
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox'],
  });
  const page = await browser.newPage();
  console.log('new tab');
  await page.goto(`${url}/${name}`, {waitUntil: 'domcontentloaded'});
  const currentUrl = page.url();
  console.log(`show page url: ${currentUrl},title: ${await page.title()}`);
  const articles = await page.evaluate(getArticles);
  const fixedUrlArticles = articles.map(article => ({
    ...omit(article, ["url"]),
    url: `${url}${article.url}`,
  }));
  console.log(`final results: ${JSON.stringify(fixedUrlArticles)}`);
  await browser.close();
  return fixedUrlArticles
};
