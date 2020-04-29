import React from 'react';
import { Route, useLocation, Redirect } from 'react-router-dom';
import { useAuth } from 'Context/auth';

const ProtectedRoute = ({ children, ...rest }) => {
  const { token } = useAuth();
  const location = useLocation();

  return (
    <Route {...rest}>
      {token ? (
        children
      ) : (
        <Redirect to={{ pathname: '/login', state: { referer: location.pathname } }} />
      )}
    </Route>
  );
};

export default ProtectedRoute;
