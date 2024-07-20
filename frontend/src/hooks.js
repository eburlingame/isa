import { useClient } from "./ClientContext";
import { useState, useEffect } from "react";

export const useClientState = () => {
  const client = useClient();
  const [clientState, setClientState] = useState(null);

  const update = (state) => {
    setClientState({ ...state });
  };

  useEffect(() => {
    if (client) {
      client.subscribeToStateChanges(update);

      return () => client.unsubscribeFromStateChanges(update);
    }
  }, [client]);

  return clientState;
};

export const useActions = () => {
  const client = useClient();

  return {
    createGame: (name) => client.createGame(name),
    joinGame: (name, gameCode) => client.joinGame(name, gameCode),
    leaveGame: () => client.leaveGame(),
    startGame: () => client.startGame(),
    endGame: () => client.endGame(),
    playCard: (cardIndex, wildColor) => client.playCard(cardIndex, wildColor),
    drawCard: () => client.drawCard(),
    doneDrawing: () => client.doneDrawing(),
  };
};
