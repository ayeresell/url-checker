import asyncio
import re
from playwright.async_api import async_playwright

async def run():
    async with async_playwright() as p:
        # Запускаем видимый браузер
        browser = await p.chromium.launch(headless=False)
        # Создаем контекст с большим таймаутом
        context = await browser.new_context()
        page = await context.new_page()

        print("🌐 Открываю страницу входа...")
        await page.goto("https://web.max.ru/-68326589202679")

        print("\n🔑 ПОЖАЛУЙСТА, ВОЙДИТЕ В АККАУНТ (отсканируйте QR-код).")
        print("После того как страница паблика полностью загрузится,")
        input("нажмите ENTER в этом терминале, чтобы начать сбор ссылок...")

        print("\n🚀 Начинаю сбор ссылок. Листаю страницу вниз...")
        
        links_found = set()
        
        # Регулярка для поиска ваших ссылок в исходном коде
        link_regex = re.compile(r'https://i\.oneme\.ru/i\?r=[a-zA-Z0-9_-]{60,75}')

        try:
            # Бесконечный скроллинг
            while True:
                # Получаем весь HTML страницы
                content = await page.content()
                new_links = link_regex.findall(content)
                
                added = 0
                for link in new_links:
                    if link not in links_found:
                        links_found.add(link)
                        with open("targets.txt", "a") as f:
                            f.write(link + "\n")
                        added += 1
                
                if added > 0:
                    print(f"📦 Найдено новых: {added} | Всего уникальных: {len(links_found)}")

                # Скроллим вниз
                await page.evaluate("window.scrollBy(0, 1000)")
                # Небольшая пауза для подгрузки контента
                await asyncio.sleep(1)
                
        except KeyboardInterrupt:
            print(f"\n🛑 Сбор остановлен. Итого собрано: {len(links_found)} уникальных ссылок.")
        finally:
            await browser.close()

if __name__ == "__main__":
    asyncio.run(run())
