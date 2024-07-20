import React, { useContext, useEffect, useRef, useState } from "react";
import GameClient from "./GameClient";

export const ClientContext = React.createContext();

export const ClientProvider = ({ url, children }) => {
  const [client, setClient] = useState(false);

  useEffect(() => {
    setClient(new GameClient("client", url));
  }, []);

  return (
    <ClientContext.Provider value={client}>{children}</ClientContext.Provider>
  );
};

export const useClient = () => {
  const client = useContext(ClientContext);

  return client;
};
