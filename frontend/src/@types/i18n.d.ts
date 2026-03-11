import ja from './core/i18n/messages/ja.json';

type Messages = typeof ja;

declare global {
  // Use type safe message keys with `next-intl`
  type IntlMessages = Messages;
}
