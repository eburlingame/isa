import { colors } from "./constants";

export const isValidPlay = (discardPileTop, wildColor, card) => {
  const topColor = getCardColorKey(discardPileTop, wildColor);
  const topNumber = getCardNumber(discardPileTop);

  if (card.startsWith("wild")) {
    return true;
  }

  if (
    (discardPileTop.endsWith("+2") && card.endsWith("+2")) ||
    (discardPileTop.endsWith("rev") && card.endsWith("rev")) ||
    (discardPileTop.endsWith("skip") && card.endsWith("skip"))
  ) {
    return true;
  }

  if (getCardColorKey(card, "") === topColor) {
    return true;
  }

  if (topNumber != -1 && getCardNumber(card) === topNumber) {
    return true;
  }

  return false;
};

export const nextValidCard = (discardPileTop, wildColor, cards, index, dir) => {
  for (let i = index + dir; i < cards.length && i >= 0; i += dir) {
    if (isValidPlay(discardPileTop, wildColor, cards[i])) {
      return i;
    }
  }

  if (dir === -1) {
    return -1;
  }

  return null;
};

export const getCardColorKey = (cardValue, wildColor) => {
  let colorKey = cardValue[0];

  if (cardValue.startsWith("wild")) {
    if (wildColor.length === 0) {
      return "R";
    }

    colorKey = wildColor;
  }

  return colorKey;
};

export const getCardNumber = (cardValue) => {
  if (
    cardValue.startsWith("wild") ||
    cardValue.endsWith("rev") ||
    cardValue.endsWith("skip") ||
    cardValue.endsWith("+2")
  ) {
    return -1;
  }

  return cardValue[1];
};

export const extractCardColor = (cardValue, wildColor) => {
  let colorKey = getCardColorKey(cardValue, wildColor);

  if (colorKey == "G") {
    return {
      color: colors.green,
      dark: colors.greenDark,
    };
  }

  if (colorKey == "R") {
    return {
      color: colors.red,
      dark: colors.redDark,
    };
  }

  if (colorKey == "B") {
    return {
      color: colors.blue,
      dark: colors.blueDark,
    };
  }

  if (colorKey == "Y") {
    return {
      color: colors.yellow,
      dark: colors.yellowDark,
    };
  }
};
