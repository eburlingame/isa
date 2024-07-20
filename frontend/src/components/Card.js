import styled from "styled-components";

import blankCardImage from "../img/cards/blank.png";

import bgBlueCardImage from "../img/cards/bg_blue.png";
import bgGreenCardImage from "../img/cards/bg_green.png";
import bgRedCardImage from "../img/cards/bg_red.png";
import bgYellowCardImage from "../img/cards/bg_yellow.png";

import number0CardImage from "../img/cards/number_0.png";
import number1CardImage from "../img/cards/number_1.png";
import number2CardImage from "../img/cards/number_2.png";
import number3CardImage from "../img/cards/number_3.png";
import number4CardImage from "../img/cards/number_4.png";
import number5CardImage from "../img/cards/number_5.png";
import number6CardImage from "../img/cards/number_6.png";
import number7CardImage from "../img/cards/number_7.png";
import number8CardImage from "../img/cards/number_8.png";
import number9CardImage from "../img/cards/number_9.png";

import pickColorCardImage from "../img/cards/pick_color.png";
import pickFourCardImage from "../img/cards/pick_four.png";
import pickTwoCardImage from "../img/cards/pick_two.png";
import reverseCardImage from "../img/cards/reverse.png";
import skipCardImage from "../img/cards/skip.png";

const boxShadow = "-0.5em 1.0em 1.0em -1.0em #000";

const CardImage = styled.img`
  object-fit: contain;
  width: 100%;
  height: 100%;

  position: relative;
  box-shadow: ${boxShadow};

  cursor: ${props => props.selectable ? "pointer" : "initial"};

  margin: 0px;
  padding: 0px;
  border-radius: 4px;
`;

const CardOverlay = styled.img`
  left: 0%;
  top: 0%;
  width: 100%;
  height: auto;

  position: absolute;

  box-shadow: ${boxShadow};
`;

const CardImageStack = styled.div`
  width: 100%;
  height: auto;

  position: relative;

  margin: 0px;
  padding: 0px;
  border-radius: 4px;
`;

const getBaseCardImage = (value) => {
  if (value === "wild") {
    return pickColorCardImage;
  }
  if (value === "wild+4") {
    return pickFourCardImage;
  }

  if (value[0] === "R") {
    return bgRedCardImage;
  }
  if (value[0] === "G") {
    return bgGreenCardImage;
  }
  if (value[0] === "B") {
    return bgBlueCardImage;
  }
  if (value[0] === "Y") {
    return bgYellowCardImage;
  }
};

const numberLayers = [
  number0CardImage,
  number1CardImage,
  number2CardImage,
  number3CardImage,
  number4CardImage,
  number5CardImage,
  number6CardImage,
  number7CardImage,
  number8CardImage,
  number9CardImage,
];

const getCardNumber = (value) => {
  const number = parseInt(value[value.length - 1]);

  if (!isNaN(number) && number >= 0 && number < numberLayers.length) {
    return number;
  }

  return null;
};

const getTopCardImage = (value) => {
  if (value.startsWith("wild")) {
    return null;
  }

  if (value.endsWith("+2")) {
    return pickTwoCardImage;
  }

  if (value.endsWith("rev")) {
    return reverseCardImage;
  }

  if (value.endsWith("skip")) {
    return skipCardImage;
  }

  // Number card
  const number = getCardNumber(value);
  if (number !== null) {
    return numberLayers[number];
  }

  return null;
};

export const Card = ({ value }) => {
  const topCardImage = getTopCardImage(value);

  return (
    <CardImageStack>
      <CardOverlay src={getBaseCardImage(value)} />
      {topCardImage && <CardOverlay src={topCardImage} />}
    </CardImageStack>
  );
};

export const BlankCard = ({ selectable }) => {
  return <CardImage selectable={selectable || false} src={blankCardImage} />;
};
