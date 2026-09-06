import { gql } from "@apollo/client";
import { useMutation } from "@apollo/client/react";
import { useState } from "react";
import { Eye, EyeOff, ArrowRight } from "lucide-react";

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

interface RegisterMutation {
  register: {
    user: {
      id: string;
      name: string;
      email: string;
      createdAt: string;
    };
    token: string;
  };
}

interface RegisterProps {
  onRegister: () => void;
  onLogin: () => void;
}

function Register({
  onRegister,
  onLogin,
}: RegisterProps) {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] =
    useState(false);

  const [register, { loading, error }] =
    useMutation<RegisterMutation>(REGISTER);

  async function handleSubmit(
    event: React.FormEvent
  ) {
    event.preventDefault();

    if (!name || !email || !password) {
      return;
    }

    try {
      const result = await register({
        variables: {
          name,
          email,
          password,
        },
      });

      const token =
        result.data?.register?.token;

      if (token) {
        localStorage.setItem(
          "stockstalk_token",
          token
        );

        onRegister();
      }
    } catch {
      // Error displayed below.
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

        <h1>Create your account</h1>

        <p className="auth-subtitle">
          Start tracking the changes that matter.
        </p>

        <form onSubmit={handleSubmit}>

          <label>Name</label>

          <input
            type="text"
            placeholder="Your name"
            value={name}
            onChange={(event) =>
              setName(event.target.value)
            }
            required
          />

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
              placeholder="Create a password"
              value={password}
              onChange={(event) =>
                setPassword(event.target.value)
              }
              required
              minLength={8}
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
              ? "Creating account..."
              : "Create account"}

            {!loading && (
              <ArrowRight size={17} />
            )}
          </button>

        </form>

        <div className="auth-switch">
          Already have an account?

          <button onClick={onLogin}>
            Sign in
          </button>
        </div>

      </div>
    </div>
  );
}

export default Register;