import puppeteer from "puppeteer"

const url = "https://zenn.dev/luvmini511"

export const start = async () => {
    console.log("start")
    const browser = await puppeteer.launch(
        {
            headless: false,
            args: ["--no-sandbox"]
        }
    )
    const page = await browser.newPage()
    console.log("new tab")
    await page.goto(url, {waitUntil: "domcontentloaded"})
    const currentUrl = page.url()
    console.log(`show page url: ${url},title: ${await page.title()}`)
    const articles = await page.evaluate(() => {
        const list = [...document.querySelectorAll(".ArticleCard_container__3qUYt")]
        return list.map(l => {
            const result = {
                title: l.querySelector(".ArticleCard_title__UnBHE")?.textContent,
                tags: [...l.querySelectorAll(".ArticleCard_topicLink__NfdwJ")].map(e => e.textContent),
                time: l.querySelector("time")?.getAttribute("datetime"),
                url: l.querySelector(".ArticleCard_mainLink__X2TOE")?.getAttribute("href")
            }
            console.log(`selected article info: ${JSON.stringify(result)}`)
            return result
        })
    })
    const fixedUrlArticles = articles.map(article => ({
        ...article, url: `${currentUrl}/${article.url}`
    }))
    console.log(`final results: ${JSON.stringify(fixedUrlArticles)}`)
    await browser.close()
}
start().then(() => console.log("end"))
