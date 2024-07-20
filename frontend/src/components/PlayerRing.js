import styled from "styled-components";
import Player from "./Player";
import arrowImage from "../img/cards/arrow.png";

const PlayerContainer = styled.div`
  position: absolute;
  left: ${(props) => props.left * 100}%;
  top: ${(props) => props.top * 100}%;
  transform: translate(-50%, -50%);
`;

const ArrowContainer = styled(PlayerContainer)``;

const ArrowImage = styled.img`
  transform: rotate(${(props) => props.rotation}rad)
    scaleX(${(props) => (props.flipped ? "-1" : "1")});
  display: ${(props) => (props.visible ? "initial" : "none")};
`;

const RADIUS_PERCENT = 0.35;

const getTablePlacementAngle = (index, count) => {
  const defaultAngle = Math.PI / 2 - (2 * Math.PI * index) / count;

  if (count === 2) {
    return (2 * Math.PI * index) / count;
  }

  return defaultAngle;
};

const getPlayerLeftPos = (index, count) =>
  0.5 + RADIUS_PERCENT * Math.cos(getTablePlacementAngle(index, count));

const getPlayerTopPos = (index, count) =>
  0.5 - RADIUS_PERCENT * Math.sin(getTablePlacementAngle(index, count));

const getArrowPlacementAngle = (index, count) =>
  getTablePlacementAngle(index, count) - Math.PI / count;

const getArrowLeftPos = (index, count) =>
  0.5 + RADIUS_PERCENT * Math.cos(getArrowPlacementAngle(index, count));

const getArrowTopPos = (index, count) =>
  0.5 - RADIUS_PERCENT * Math.sin(getArrowPlacementAngle(index, count));

const arrowRotationAngle = (direction, index, count) => {
  const angle = -getArrowPlacementAngle(index, count);

  if (direction > 0) {
    return angle - 0.2;
  } else {
    return angle + Math.PI + 0.2;
  }
};

const nextPlayer = (index, dir, count) => {
  const next = index + dir;
  if (next >= count) {
    return next - count;
  }
  if (next < 0) {
    return next + count;
  }
  return next;
};

export default ({ activePlayer, direction, players, discardPile }) => {
  return (
    <>
      {players.map(({ name, numCards }, index) => (
        <PlayerContainer
          key={name}
          left={getPlayerLeftPos(index, players.length)}
          top={getPlayerTopPos(index, players.length)}
        >
          <Player
            playerActive={activePlayer === index}
            name={name}
            numCards={numCards}
          />
        </PlayerContainer>
      ))}

      {players.map(({ name }, index) => (
        <ArrowContainer
          key={name}
          left={getArrowLeftPos(index, players.length)}
          top={getArrowTopPos(index, players.length)}
        >
          <ArrowImage
            visible={
              direction === 1
                ? activePlayer === index
                : nextPlayer(activePlayer, direction, players.length) == index
            }
            src={arrowImage}
            flipped={direction < 0}
            rotation={arrowRotationAngle(direction, index, players.length)}
          />
        </ArrowContainer>
      ))}
    </>
  );
};
