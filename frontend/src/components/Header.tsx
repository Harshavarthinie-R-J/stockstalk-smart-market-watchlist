import {
  User,
  LogOut,
  ChevronDown,
  X,
  Search,
  Sun,
  Moon,
  Bell,
  CheckCheck,
} from "lucide-react";

import {
  useEffect,
  useRef,
  useState,
} from "react";

interface HeaderProps {
  onLogout: () => void;
}

interface StoredUser {
  id?: string;
  name?: string;
  email?: string;
}

export default function Header({
  onLogout,
}: HeaderProps) {
  const [menuOpen, setMenuOpen] =
    useState(false);

  const [profileOpen, setProfileOpen] =
    useState(false);

  const [notificationsOpen, setNotificationsOpen] =
    useState(false);

  const [search, setSearch] =
    useState("");

  const [darkMode, setDarkMode] =
    useState(true);

  const profileRef =
    useRef<HTMLDivElement>(null);

  const notificationRef =
    useRef<HTMLDivElement>(null);

  const [user, setUser] =
    useState<StoredUser>({
      name: "Investor",
      email: "",
    });


  /* =====================================================
     LOAD USER
     ===================================================== */

  useEffect(() => {
    const storedUser =
      localStorage.getItem(
        "stockstalk_user"
      );

    if (!storedUser) {
      return;
    }

    try {
      const parsed =
        JSON.parse(storedUser);

      setUser({
        id: parsed.id,
        name:
          parsed.name ||
          parsed.username ||
          "Investor",
        email:
          parsed.email || "",
      });
    } catch (error) {
      console.error(
        "Unable to load user:",
        error
      );
    }
  }, []);


  /* =====================================================
     LOAD THEME
     ===================================================== */

  useEffect(() => {
    const savedTheme =
      localStorage.getItem(
        "stockstalk_theme"
      );

    const isDark =
      savedTheme !== "light";

    setDarkMode(isDark);

    document.documentElement.setAttribute(
      "data-theme",
      isDark ? "dark" : "light"
    );
  }, []);


  /* =====================================================
     CLOSE DROPDOWNS WHEN CLICKING OUTSIDE
     ===================================================== */

  useEffect(() => {
    function handleOutsideClick(
      event: MouseEvent
    ) {
      const target =
        event.target as Node;

      if (
        profileRef.current &&
        !profileRef.current.contains(
          target
        )
      ) {
        setMenuOpen(false);
      }

      if (
        notificationRef.current &&
        !notificationRef.current.contains(
          target
        )
      ) {
        setNotificationsOpen(false);
      }
    }

    document.addEventListener(
      "mousedown",
      handleOutsideClick
    );

    return () => {
      document.removeEventListener(
        "mousedown",
        handleOutsideClick
      );
    };
  }, []);


  /* =====================================================
     THEME
     ===================================================== */

  function toggleTheme() {
    const nextDark =
      !darkMode;

    setDarkMode(nextDark);

    document.documentElement.setAttribute(
      "data-theme",
      nextDark ? "dark" : "light"
    );

    localStorage.setItem(
      "stockstalk_theme",
      nextDark
        ? "dark"
        : "light"
    );
  }


  /* =====================================================
     LOGOUT
     ===================================================== */

  function handleLogout(
    event?: React.MouseEvent
  ) {
    event?.preventDefault();
    event?.stopPropagation();

    setMenuOpen(false);
    setProfileOpen(false);
    setNotificationsOpen(false);
    localStorage.removeItem(
      "stockstalk_token"
    );

    localStorage.removeItem(
      "stockstalk_user"
    );

    onLogout();
  }


  /* =====================================================
     SEARCH
     ===================================================== */

  function handleSearchKeyDown(
    event: React.KeyboardEvent<HTMLInputElement>
  ) {
    if (event.key === "Enter") {
      const value =
        search.trim();

      if (value) {
        console.log(
          "Search:",
          value
        );
      }
    }
  }


  return (
    <>
      <header className="stockstalk-header">

        {/* =================================================
            LEFT - BRAND
            ================================================= */}

        <div className="stockstalk-brand">

          <div className="stockstalk-logo">
            <img
              src="/groww-logo.png"
              alt="Groww"
            />
          </div>

          <div className="stockstalk-brand-text">

            <strong>
              StockStalk
            </strong>

            <span>
              Your stocks. What changed. What matters.
            </span>

          </div>

        </div>


        {/* =================================================
            CENTER - SEARCH
            ================================================= */}

        <div className="stockstalk-header-search">

          <div className="header-search-box">

            <Search
              size={17}
              className="header-search-icon"
            />

            <input
              type="text"
              value={search}
              onChange={(event) =>
                setSearch(
                  event.target.value
                )
              }
              onKeyDown={
                handleSearchKeyDown
              }
              placeholder="Search stocks, companies..."
              aria-label="Search stocks"
            />

            {search && (
              <button
                type="button"
                className="header-search-clear"
                onClick={() =>
                  setSearch("")
                }
              >
                <X size={14} />
              </button>
            )}

          </div>

        </div>


        {/* =================================================
            RIGHT ACTIONS
            ================================================= */}

        <div className="stockstalk-header-actions">


          {/* =============================================
              THEME BUTTON
              ============================================= */}

          <button
            type="button"
            className="header-icon-button"
            onClick={toggleTheme}
            title={
              darkMode
                ? "Switch to light mode"
                : "Switch to dark mode"
            }
            aria-label="Toggle theme"
          >

            {darkMode ? (
              <Sun size={18} />
            ) : (
              <Moon size={18} />
            )}

          </button>


          {/* =============================================
              NOTIFICATION
              ============================================= */}

          <div
            className="header-notification-wrapper"
            ref={notificationRef}
          >

            <button
              type="button"
              className={`header-icon-button ${
                notificationsOpen
                  ? "header-icon-button-active"
                  : ""
              }`}
              onClick={() =>
                setNotificationsOpen(
                  (previous) =>
                    !previous
                )
              }
              title="Notifications"
              aria-label="Notifications"
            >

              <Bell size={18} />

              <span className="notification-dot" />

            </button>


            {notificationsOpen && (
              <div className="notification-menu">

                <div className="notification-menu-header">

                  <div>
                    <strong>
                      Notifications
                    </strong>

                    <span>
                      Market activity
                    </span>
                  </div>

                  <CheckCheck
                    size={16}
                  />

                </div>


                <div className="notification-divider" />


                <div className="notification-empty">

                  <div className="notification-empty-icon">
                    <Bell size={20} />
                  </div>

                  <strong>
                    You're all caught up
                  </strong>

                  <span>
                    Important watchlist events
                    will appear here.
                  </span>

                </div>

              </div>
            )}

          </div>


          {/* =============================================
              INVESTOR PROFILE
              ============================================= */}

          <div
            className="profile-wrapper"
            ref={profileRef}
          >

            <button
              type="button"
              className={`profile-trigger ${
                menuOpen
                  ? "profile-trigger-active"
                  : ""
              }`}
              onClick={(event) => {
                event.stopPropagation();

                setMenuOpen(
                  (previous) =>
                    !previous
                );

                setNotificationsOpen(
                  false
                );
              }}
              aria-expanded={menuOpen}
            >

              <span className="profile-avatar">
                <User size={19} />
              </span>

              <span className="profile-trigger-name">
                {user.name ||
                  "Investor"}
              </span>

              <ChevronDown
                size={15}
                className={
                  menuOpen
                    ? "profile-chevron profile-chevron-open"
                    : "profile-chevron"
                }
              />

            </button>


            {/* =========================================
                PROFILE DROPDOWN
                ========================================= */}

            {menuOpen && (
              <div
                className="profile-menu"
                onClick={(event) =>
                  event.stopPropagation()
                }
              >

                <div className="profile-menu-user">

                  <div className="profile-menu-avatar">
                    <User size={19} />
                  </div>

                  <div className="profile-menu-user-info">

                    <strong>
                      {user.name ||
                        "Investor"}
                    </strong>

                    <span>
                      {user.email ||
                        "Logged in user"}
                    </span>

                  </div>

                </div>


                <div className="profile-menu-divider" />


                {/* PROFILE */}

                <button
                  type="button"
                  className="profile-menu-item"
                  onClick={() => {
                    setMenuOpen(false);
                    setProfileOpen(true);
                  }}
                >

                  <User size={17} />

                  <span>
                    View Profile
                  </span>

                </button>


                {/* LOGOUT */}

                <button
                  type="button"
                  className="profile-menu-item profile-logout"
                  
                  onClick={(event) => {
                    event.preventDefault();
                    event.stopPropagation();

                    localStorage.removeItem("stockstalk_token");
                    localStorage.removeItem("stockstalk_user");

                    window.location.href = "/";
                  }}
                >

                  <LogOut size={17} />

                  <span>
                    Logout
                  </span>

                </button>

              </div>
            )}

          </div>

        </div>

      </header>


      {/* =================================================
          PROFILE MODAL
          ================================================= */}

      {profileOpen && (
        <div
          className="profile-modal-overlay"
          onClick={() =>
            setProfileOpen(false)
          }
        >

          <div
            className="profile-modal"
            onClick={(event) =>
              event.stopPropagation()
            }
          >

            <div className="profile-modal-header">

              <div>

                <div className="section-kicker">
                  ACCOUNT
                </div>

                <h2>
                  Profile
                </h2>

              </div>

              <button
                type="button"
                className="profile-modal-close"
                onClick={() =>
                  setProfileOpen(false)
                }
              >
                <X size={18} />
              </button>

            </div>


            <div className="profile-modal-user">

              <div className="profile-large-avatar">
                <User size={30} />
              </div>

              <div>

                <h3>
                  {user.name ||
                    "Investor"}
                </h3>

                <p>
                  StockStalk Investor
                </p>

              </div>

            </div>


            <div className="profile-details">

              <div className="profile-detail">

                <span>
                  Name
                </span>

                <strong>
                  {user.name ||
                    "Investor"}
                </strong>

              </div>


              <div className="profile-detail">

                <span>
                  Email
                </span>

                <strong>
                  {user.email ||
                    "Not available"}
                </strong>

              </div>

            </div>


            <button
              type="button"
              className="profile-modal-logout"
              onClick={handleLogout}
            >

              <LogOut size={17} />

              Logout

            </button>

          </div>

        </div>
      )}

    </>
  );
}