import {
  ApolloClient,
  InMemoryCache,
  HttpLink,
} from "@apollo/client";

import {
  setContext,
} from "@apollo/client/link/context";

const httpLink = new HttpLink({
  uri: import.meta.env.VITE_GRAPHQL_URL ||
    "/graphql",
});

const authLink = setContext(
  (_, { headers }) => {
    const token =
      localStorage.getItem(
        "stockstalk_token"
      );

    return {
      headers: {
        ...headers,

        ...(token
          ? {
              Authorization:
                `Bearer ${token}`,
            }
          : {}),
      },
    };
  }
);

const client = new ApolloClient({
  link: authLink.concat(httpLink),

  cache: new InMemoryCache(),
});

export default client;