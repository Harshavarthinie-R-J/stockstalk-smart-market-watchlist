import { gql } from "@apollo/client";
import { useMutation } from "@apollo/client/react";
import { Loader2, LogIn, UserPlus } from "lucide-react";
import { useState } from "react";
import useDayNightTheme from "./hooks/useDayNightTheme";
import Dashboard from "./pages/Dashboard";

const LOGIN = gql`
  mutation Login(
    $email: String!
    $password: String!
  ) {
    login(
      email: $email
      password: $password
    ) {
      token
      user {
        id
        name
        email
      }
    }
  }
`;

const REGISTER = gql`
  mutation Register(
    $name: String!
    $email: String!
    $password: String!
  ) {
    register(
      name: $name
      email: $email
      password: $password
    ) {
      token
      user {
        id
        name
        email
      }
    }
  }
`;

interface AuthResponse {
  token: string;
  user?: {
    id: string;
    name: string;
    email: string;
  } | null;
}

interface LoginData {
  login: AuthResponse;
}

interface RegisterData {
  register: AuthResponse;
}

function AuthPage({
  onLogin,
}: {
  onLogin: () => void;
}) {
  const [isRegister, setIsRegister] =
    useState(false);

  const [name, setName] = useState("");

  const [email, setEmail] =
    useState("");

  const [password, setPassword] =
    useState("");

  const [errorMessage, setErrorMessage] =
    useState("");

  const [login, { loading: loginLoading }] =
    useMutation<
      LoginData,
      {
        email: string;
        password: string;
      }
    >(LOGIN);

  const [
    register,
    { loading: registerLoading },
  ] = useMutation<
    RegisterData,
    {
      name: string;
      email: string;
      password: string;
    }
  >(REGISTER);

  const loading =
    loginLoading || registerLoading;

  async function handleSubmit(
    event: React.FormEvent
  ) {
    event.preventDefault();

    setErrorMessage("");

    const cleanEmail =
      email.trim();

    const cleanName =
      name.trim();

    if (!cleanEmail || !password) {
      setErrorMessage(
        "Please enter your email and password."
      );
      return;
    }

    if (
      isRegister &&
      !cleanName
    ) {
      setErrorMessage(
        "Please enter your name."
      );
      return;
    }

    try {
      if (isRegister) {
        const result =
          await register({
            variables: {
              name: cleanName,
              email: cleanEmail,
              password,
            },
          });

        const auth =
          result.data?.register;

        if (!auth?.token) {
          throw new Error(
            "Registration did not return a token."
          );
        }

        localStorage.setItem(
          "stockstalk_token",
          auth.token
        );

        if (auth.user) {
          localStorage.setItem(
            "stockstalk_user",
            JSON.stringify(auth.user)
          );
        }
      } else {
        const result =
          await login({
            variables: {
              email: cleanEmail,
              password,
            },
          });

        const auth =
          result.data?.login;

        if (!auth?.token) {
          throw new Error(
            "Login did not return a token."
          );
        }

        localStorage.setItem(
          "stockstalk_token",
          auth.token
        );

        if (auth.user) {
          localStorage.setItem(
            "stockstalk_user",
            JSON.stringify(auth.user)
          );
        }
      }

      onLogin();
    } catch (error) {
      console.error(
        "Authentication error:",
        error
      );

      const message =
        error instanceof Error
          ? error.message
          : "Unable to authenticate.";

      setErrorMessage(message);
    }
  }

  return (
    <div className="auth-page">
      <div className="auth-card">

        <div className="auth-logo">
          S
        </div>

        <div className="auth-brand">
          StockStalk
        </div>

        <div className="auth-tagline">
          Your stocks. What changed.
          What matters.
        </div>

        <div className="auth-heading">
          <h1>
            {isRegister
              ? "Create your account"
              : "Welcome back"}
          </h1>

          <p>
            {isRegister
              ? "Start building your smart market watchlist."
              : "Sign in to continue to your watchlists."}
          </p>
        </div>

        <form
          className="auth-form"
          onSubmit={handleSubmit}
        >

          {isRegister && (
            <div className="auth-field">
              <label htmlFor="name">
                Name
              </label>

              <input
                id="name"
                type="text"
                value={name}
                placeholder="Your name"
                disabled={loading}
                onChange={(event) =>
                  setName(event.target.value)
                }
              />
            </div>
          )}

          <div className="auth-field">
            <label htmlFor="email">
              Email
            </label>

            <input
              id="email"
              type="email"
              value={email}
              placeholder="you@example.com"
              disabled={loading}
              onChange={(event) =>
                setEmail(event.target.value)
              }
            />
          </div>

          <div className="auth-field">
            <label htmlFor="password">
              Password
            </label>

            <input
              id="password"
              type="password"
              value={password}
              placeholder="Enter your password"
              disabled={loading}
              onChange={(event) =>
                setPassword(event.target.value)
              }
            />
          </div>

          {errorMessage && (
            <div className="auth-error">
              {errorMessage}
            </div>
          )}

          <button
            type="submit"
            className="auth-submit-button"
            disabled={loading}
          >
            {loading ? (
              <>
                <Loader2
                  size={17}
                  className="spin"
                />
                Please wait...
              </>
            ) : isRegister ? (
              <>
                <UserPlus size={17} />
                Create account
              </>
            ) : (
              <>
                <LogIn size={17} />
                Sign in
              </>
            )}
          </button>

        </form>

        <div className="auth-switch">
          <span>
            {isRegister
              ? "Already have an account?"
              : "Don't have an account?"}
          </span>

          <button
            type="button"
            onClick={() => {
              setIsRegister(
                (current) => !current
              );

              setErrorMessage("");
            }}
          >
            {isRegister
              ? "Sign in"
              : "Create account"}
          </button>
        </div>

      </div>
    </div>
  );
}

function App() {
  useDayNightTheme();

  const [loggedIn, setLoggedIn] =
    useState(
      Boolean(
        localStorage.getItem(
          "stockstalk_token"
        )
      )
    );

  function handleLogout() {
    localStorage.removeItem(
      "stockstalk_token"
    );

    localStorage.removeItem(
      "stockstalk_user"
    );

    setIsAuthenticated(false);
  }

  function handleLogin() {
    setLoggedIn(true);
  }

  if (!loggedIn) {
    return (
      <AuthPage
        onLogin={handleLogin}
      />
    );
  }

  return (
    <Dashboard
      onLogout={handleLogout}
    />
  );
}

export default App;

function setIsAuthenticated(arg0: boolean) {
  throw new Error("Function not implemented.");
}
