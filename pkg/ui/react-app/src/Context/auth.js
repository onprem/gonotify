import { createContext, useContext } from 'react';

export const AuthContext = createContext();

export const AuthProvider = AuthContext.Provider;

export function useAuth() {
  return useContext(AuthContext);
}
