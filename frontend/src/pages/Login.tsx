import { gql } from "@apollo/client";
import { useMutation } from "@apollo/client/react";
import { useState } from "react";
import { Eye, EyeOff, ArrowRight } from "lucide-react";

const LOGIN = gql`
  mutation Login($email: String!, $password: String!) {
    login(email: $email, password: $password) {
      user {
        id
        name
        email
        createdAt
      }
      token
    }
  }
`;

interface LoginProps {
  onLogin: () => void;
  onRegister: () => void;
}

interface LoginData {
  login: {
    user: {
      id: string;
      name: string;
      email: string;
      createdAt: string;
    };
    token: string;
  };
}

interface LoginVariables {
  email: string;
  password: string;
}

function Login({ onLogin, onRegister }: LoginProps) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);

  const [login, { loading, error }] = useMutation<
    LoginData,
    LoginVariables
  >(LOGIN);

  async function handleSubmit(
    event: React.FormEvent
  ) {
    event.preventDefault();

    if (!email || !password) {
      return;
    }

    try {
      const result = await login({
        variables: {
          email,
          password,
        },
      });

      const token = result.data?.login?.token;

      if (token) {
        localStorage.setItem(
          "stockstalk_token",
          token
        );

        onLogin();
      }
    } catch {
      // Error is displayed below the form.
    }
  }

  return (
    <div className="auth-page">
      <div className="auth-card">

        <img
          className="login-groww-logo"
          src="/groww-logo.png"
          alt="Groww"
        />      

        <h1>Welcome back</h1>

        <p className="auth-subtitle">
          Sign in to see what changed in your
          watchlists.
        </p>

        <form onSubmit={handleSubmit}>

          <label>Email</label>

          <input
            type="email"
            placeholder="you@example.com"
            value={email}
            onChange={(event) =>
              setEmail(event.target.value)
            }
            required
          />

          <label>Password</label>

          <div className="password-input">

            <input
              type={
                showPassword
                  ? "text"
                  : "password"
              }
              placeholder="Enter your password"
              value={password}
              onChange={(event) =>
                setPassword(event.target.value)
              }
              required
            />

            <button
              type="button"
              onClick={() =>
                setShowPassword(!showPassword)
              }
            >
              {showPassword ? (
                <EyeOff size={18} />
              ) : (
                <Eye size={18} />
              )}
            </button>

          </div>

          {error && (
            <div className="auth-error">
              {error.message}
            </div>
          )}

          <button
            className="auth-submit"
            type="submit"
            disabled={loading}
          >
            {loading
              ? "Signing in..."
              : "Sign in"}

            {!loading && (
              <ArrowRight size={17} />
            )}
          </button>

        </form>

        <div className="auth-switch">
          Don't have an account?

          <button onClick={onRegister}>
            Create one
          </button>
        </div>

      </div>
    </div>
  );
}

export default Login;