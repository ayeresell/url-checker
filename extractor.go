import re

# Регулярное выражение для поиска ваших ссылок
# Ищем https://i.oneme.ru/i?r= и затем 67 символов Base64
link_pattern = r'https://i\.oneme\.ru/i\?r=[a-zA-Z0-9_-]{60,75}'

def extract_links(input_file, output_file):
    try:
        with open(input_file, 'r', encoding='utf-8') as f:
            content = f.read()
        
        # Находим все совпадения
        links = re.findall(link_pattern, content)
        
        # Оставляем только уникальные
        unique_links = list(set(links))
        
        with open(output_file, 'w', encoding='utf-8') as f:
            for link in unique_links:
                f.write(link + '\n')
        
        print(f"✅ Успех! Найдено уникальных ссылок: {len(unique_links)}")
        print(f"📁 Результаты сохранены в: {output_file}")
        
    except FileNotFoundError:
        print(f"❌ Ошибка: Файл '{input_file}' не найден. Создайте его и положите туда текст.")

if __name__ == "__main__":
    # Инструкция:
    # 1. Создайте файл input.txt в этой папке
    # 2. Вставьте туда весь текст, где есть ссылки
    # 3. Запустите этот скрипт
    extract_links('input.txt', 'targets.txt')
