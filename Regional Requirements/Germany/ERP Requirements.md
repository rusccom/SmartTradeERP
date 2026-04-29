# Германия: технические требования к ERP/учетной системе

Документ описывает ориентиры для развития ядра SmartTrade ERP под Германию. Это не юридическое заключение и не активный backlog: перед реализацией требования нужно подтверждать с налоговым консультантом, юристом и актуальными официальными источниками.

## Текущий вывод по ядру

Текущая архитектура подходит как основа для внутреннего складского и операционного учета: товары, варианты, склады, документы, продажи, возвраты, смены, платежи, cost ledger, tenant isolation.

Для Германии на это ядро можно опираться только как на стартовую техническую базу. Для фискального учета, POS, B2B-инвойсинга и e-commerce compliance нужны отдельные слои: VAT, неизменяемый audit trail, invoice/e-invoice, TSE, экспорт для налогового консультанта и полноценная модель юридических документов.

## 1. VAT / Umsatzsteuer

Нужно добавить налоговую модель на уровне tenant, товара, строки документа и итогов документа.

Минимальная модель:
- `currency`, по умолчанию `EUR`;
- `tax_rate`: 19%, 7%, 0%, exempt, reverse charge, intra-EU supply;
- `net_amount`, `tax_amount`, `gross_amount`;
- правило округления по строке и по документу;
- налоговая категория товара;
- VAT ID / USt-IdNr для tenant и B2B-клиентов;
- признак B2B/B2C клиента;
- страна клиента и место поставки;
- OSS/IOSS flags для cross-border B2C;
- snapshot налоговых параметров в документе, чтобы изменения справочников не меняли историю.

Текущий пробел: сейчас в ядре есть `unit_price` и `total_amount`, но нет net/gross/VAT split и налоговых ставок.

## 2. Неизменяемость и GoBD

Для Германии критично, чтобы учетные записи были полными, прослеживаемыми, своевременными, упорядоченными и неизменяемыми.

Технические требования:
- posted-документы нельзя физически переписывать как единственную версию истины;
- исправления проводить через correction/reversal documents;
- хранить immutable audit log: кто, когда, что изменил, причина, old/new snapshot;
- ledger должен быть append-only либо иметь отдельный журнал корректировок;
- удаление бизнес-документов после posting запретить;
- закрытые периоды блокировать от изменения;
- все номера документов должны быть последовательными и объяснимыми;
- экспорт данных должен быть воспроизводимым.

Текущий риск: retro-update posted-документов и удаление ledger rows удобны для управленческого учета, но плохо подходят как фискальная история.

## 3. POS, смены, касса, TSE

Если система используется как касса или POS в Германии, нужен отдельный фискальный слой.

Минимальные требования:
- интеграция с сертифицированной TSE;
- старт/завершение TSE transaction на каждую кассовую операцию;
- хранение TSE transaction number, signature, start/end time, serial number;
- обязательные поля кассового чека;
- печать/электронная выдача чека;
- обработка отказа TSE без потери аудита;
- DSFinV-K/export для проверок;
- запрет тихого изменения кассовых операций после фискализации.

Если SmartTrade остается только back-office ERP без кассовой функции, TSE может быть вне ядра, но граница интеграции с POS должна быть явно описана.

## 4. Invoices и E-Rechnung

Для B2B в Германии нужно проектировать invoice subsystem отдельно от складских документов.

Минимальные требования:
- invoice, credit note, cancellation invoice;
- seller legal data snapshot;
- buyer legal data snapshot;
- invoice issue date, supply date/period, due date;
- VAT breakdown по ставкам;
- invoice number sequence;
- payment terms;
- поддержка XRechnung и/или ZUGFeRD/Factur-X;
- прием входящих E-Rechnung и хранение исходного XML/PDF;
- validation status и технический audit trail;
- связь invoice с shipment/order/payment/document.

Важно: складской `SALE` документ не должен автоматически считаться юридической Rechnung без invoice layer.

## 5. Бухгалтерский слой

Cost ledger полезен для себестоимости, но не заменяет бухгалтерию.

Для немецкого рынка потребуются:
- chart of accounts, например SKR03/SKR04 как настраиваемый mapping;
- journal entries;
- AR/AP;
- supplier invoices;
- bank transactions и reconciliation;
- VAT return data;
- DATEV export;
- period closing;
- audit-safe corrections.

Практичный путь: сначала сделать экспорт для Steuerberater/DATEV, а не строить полную бухгалтерию сразу.

## 6. E-commerce слой

Если система будет обслуживать интернет-магазин, нужно поддержать требования до момента заказа.

Технические элементы:
- cart/checkout/order model;
- отображение gross prices для B2C;
- shipping cost, delivery restrictions, accepted payment methods до заказа;
- кнопка заказа с немецкой формулировкой вроде `zahlungspflichtig bestellen`;
- withdrawal/cancellation flow;
- order confirmation;
- AGB, Widerrufsbelehrung, Datenschutz, Impressum;
- cookie consent для нестрого необходимых cookies;
- product safety data для B2C-товаров.

Сейчас в проекте есть landing и placeholders для legal links, но нет полноценного checkout/order flow.

## 7. Product safety и товарные данные

Для B2C-торговли в ЕС нужно учитывать product safety requirements.

Для каталога стоит предусмотреть:
- manufacturer / responsible person;
- address/contact responsible economic operator;
- warnings and safety instructions;
- product identifiers: SKU, GTIN/EAN, batch/lot/serial;
- country of origin при необходимости;
- attachments: manuals, safety docs, certificates;
- recall/withdrawal workflow;
- связь партии товара с поставкой и продажей.

## 8. GDPR / Datenschutz

Система хранит персональные данные клиентов и пользователей, поэтому нужен отдельный privacy/data governance слой.

Минимальные технические требования:
- data inventory по категориям персональных данных;
- purpose/legal basis для обработки;
- retention policies;
- export data subject data;
- delete/anonymize where legally possible;
- access logs для чувствительных действий;
- role-based access control по действиям, а не только роль в токене;
- MFA для admin/owner;
- короткие access tokens, refresh/session model;
- secure cookie strategy вместо long-lived localStorage tokens для production;
- tenant-level DPA/subprocessor documentation вне кода.

Текущие риски: долгий JWT TTL, хранение токенов в `localStorage`, широкий CORS, роли почти не применяются в маршрутах.

## 9. Архитектурный план внедрения

Рекомендуемый порядок:

1. Добавить country profile: `DE`, currency, locale, tax defaults, invoice numbering rules.
2. Ввести VAT engine и snapshot налогов в документах.
3. Разделить warehouse documents, orders, invoices, payments и ledger postings.
4. Перевести posted-документы на correction/reversal model.
5. Добавить immutable audit log.
6. Добавить invoice subsystem и E-Rechnung export/import.
7. Добавить DATEV/export слой.
8. Описать POS boundary: либо TSE integration, либо явный non-POS режим.
9. Добавить e-commerce checkout compliance только если система реально становится интернет-магазином.
10. Усилить security/GDPR слой.

## 10. Что не делать в ядре преждевременно

- Не смешивать складской `SALE` с юридической Rechnung.
- Не пришивать TSE внутрь всех документов, если POS может быть отдельным модулем.
- Не делать VAT как одно поле `tax_percent` без snapshots и налоговых категорий.
- Не позволять редактировать posted-документ без correction/audit модели.
- Не строить полную бухгалтерию до понимания DATEV/Steuerberater workflow.

## Официальные ориентиры

- UStG § 12: стандартная и сниженная ставки Umsatzsteuer.
- BMF FAQ E-Rechnung: обязательная E-Rechnung для B2B с 1 января 2025.
- BMF GoBD: требования к электронным книгам, записям, документам и доступу к данным.
- KassenSichV / AO § 146a: требования к электронным кассовым системам и TSE.
- BGB § 312j: обязанности в e-commerce checkout и платежная кнопка.
- DDG § 5: Impressum / allgemeine Informationspflichten.
- PAngV: правила отображения цен, VAT и дополнительных расходов.
- TDDDG § 25 и GDPR: cookies, consent, privacy/security.
- EU OSS: VAT reporting для cross-border B2C продаж внутри ЕС.
- Regulation (EU) 2023/988: product safety для consumer goods, включая online sales.

Проверенные источники на дату создания документа:
- https://www.gesetze-im-internet.de/ustg_1980/__12.html
- https://www.bundesfinanzministerium.de/Content/DE/FAQ/e-rechnung.html
- https://www.bundesfinanzministerium.de/Content/DE/Downloads/BMF_Schreiben/Weitere_Steuerthemen/Abgabenordnung/AO-Anwendungserlass/2024-03-11-aenderung-gobd.html
- https://www.bundesfinanzministerium.de/Content/DE/Gesetzestexte/Gesetze_Gesetzesvorhaben/Abteilungen/Abteilung_IV/18_Legislaturperiode/Gesetze_Verordnungen/2017-10-06-KassenSichV/0-Verordnung.html
- https://www.gesetze-im-internet.de/bgb/__312j.html
- https://www.gesetze-im-internet.de/ddg/__5.html
- https://www.gesetze-im-internet.de/pangv_2022/BJNR492110021.html
- https://www.gesetze-im-internet.de/ttdsg/__25.html
- https://europa.eu/youreurope/business/taxation/vat/one-stop-shop/index_de.htm
- https://eur-lex.europa.eu/EN/legal-content/summary/general-product-safety-regulation-2023.html
