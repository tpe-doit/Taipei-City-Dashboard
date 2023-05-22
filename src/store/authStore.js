import { defineStore } from "pinia";
import axios from "axios";
import { md5 } from "../assets/utilityFunctions/md5";
import router from "../router";
import { useContentStore } from "./contentStore";

export const useAuthStore = defineStore("auth", {
  state: () => ({
    // cmd f userDocs
    user: null,
    auth: 0,
    tokens: {
      access_token: localStorage.getItem("access_token") || "",
      refresh_token: localStorage.getItem("refresh_token") || "",
      token_type: localStorage.getItem("token_type") || "",
    },
    errorMessage: "",
  }),
  getters: {
    userEssentials() {
      const accountType = ["Email 用戶", "台北通用戶", "Taipei On 用戶"];
      return {
        name: this.user.name,
        type: accountType[this.user.type],
      };
    },
  },
  actions: {
    // Sends email and password to backend to request tokens. Success -> setTokens / setUser / reroute to dashboard ; Fail -> Display Error
    handleEmailLogin(email, password) {
      axios
        .post(
          `/api_server/access/token`,
          {},
          {
            headers: {
              Authorization: `Basic ${btoa(`${email}:${md5(password)}`)}`,
            },
          }
        )
        .then((rs) => {
          this.setTokens(
            rs.data.access_token,
            rs.data.refresh_token,
            rs.data.token_type
          );
          this.setUser();
          router.replace({ path: "/Dashboard" });
        })
        .catch((error) => {
          console.log(error);
        });
    },
    // Sends request to backend to logout. Success or Fail -> Clear localStorage / Clear Store / Reroute to login
    handleLogout() {
      if (this.tokens.access_token) {
        axios
          .post(`/api_server/access/logout`, {
            headers: {
              Authorization: `${this.tokens.token_type} ${this.tokens.access_token}`,
            },
          })
          .then((rs) => {
            console.log(rs);
          })
          .catch((error) => {
            console.log(error);
          });
      }
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
      localStorage.removeItem("token_type");
      this.executeClearStore();
      router.replace({ path: "/login" });
    },
    // Posts the refresh token to the server to request a new pair of tokens then sets user info
    executeRefreshTokens() {
      axios({
        method: "post",
        url: "/api_server/access/refresh",
        headers: {
          Authorization: `${this.tokens.token_type} ${this.tokens.refresh_token}`,
        },
      })
        .then((rs) => {
          const newToken = rs.data;
          this.setTokens(
            newToken.accessToken,
            newToken.refresh_token,
            newToken.token_type
          );
          axios.defaults.headers.common[
            "Authorization"
          ] = `${newToken.token_type} ${newToken.access_token}`;
          axios
            .get(`/api_server/manager/authuser`)
            .then((rs) => {
              this.user = rs.data;
              this.auth = rs.data.gid == 1 ? 2 : 1;
            })
            .catch((e) => {
              this.handleLogout();
            });
        })
        .catch((e) => {
          this.handleLogout();
        });
    },
    // clears the entire store
    executeClearStore() {
      this.user = null;
      this.auth = 0;
      this.tokens = {
        access_token: localStorage.getItem("access_token") || "",
        refresh_token: localStorage.getItem("refresh_token") || "",
        token_type: localStorage.getItem("token_type") || "",
      };
      this.dashboards = {
        fixed: [],
        customized: [],
      };
    },
    // sets the tokens into the store and local storage
    setTokens(access_token, refresh_token, token_type) {
      this.tokens = {
        access_token: access_token,
        refresh_token: refresh_token,
        token_type: token_type,
      };
      localStorage.setItem("access_token", access_token);
      localStorage.setItem("refresh_token", refresh_token);
      localStorage.setItem("token_type", token_type);
    },
    // requests user info from the server and adds it to the store
    setUser() {
      const contentStore = useContentStore();
      if (!this.tokens.access_token || !this.tokens.refresh_token) {
        this.handleLogout();
        return;
      }
      axios.defaults.headers.common[
        "Authorization"
      ] = `${this.tokens.token_type} ${this.tokens.access_token}`;
      axios
        .get(`/api_server/manager/authuser`)
        .then((rs) => {
          this.user = rs.data;
          this.auth = rs.data.gid == 1 ? 2 : 1;
        })
        .catch((e) => {
          console.log(e);
          this.executeRefreshTokens();
        });
    },
  },
});

// userDocs:
/*
{
  activated_at: Date,
  created_at: Date,
  email: String,
  gid: Number,
  id: Number,
  is_blacklist: Boolean,
  is_whitelist: Boolean,
  login_at: Date,
  name: String,
  status: Number,
  type: Number,
  updated_at: Date,
  updated_by: Date,
  taipei_pass: Object,
}
*/
