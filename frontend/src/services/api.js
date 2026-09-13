const API_BASE = "http://localhost:8080";

async function request(url, options = {}) {
  const response = await fetch(`${API_BASE}${url}`, {
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    ...options,
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || "Request failed");
  }

  return response.json();
}

export const participants = {
  getAll: () => request("/participants"),
  getById: (id) => request(`/participants/${id}`),
  create: (data) =>
    request("/participants", { method: "POST", body: JSON.stringify(data) }),
  update: (id, data) =>
    request(`/participants/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),
  activate: (id) => request(`/participants/${id}/activate`, { method: "PUT" }),
  delete: (id) => request(`/participants/${id}`, { method: "DELETE" }),
};

export const contributions = {
  getAll: () => request("/contributions"),
  getById: (id) => request(`/contributions/${id}`),
  create: (data) =>
    request("/contributions", { method: "POST", body: JSON.stringify(data) }),
  update: (id, data) =>
    request(`/contributions/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),
  delete: (id) => request(`/contributions/${id}`, { method: "DELETE" }),
  getByParticipant: (id) => request(`/contributions/participant/${id}`),
  getByPeriod: (month, year) =>
    request(`/contributions/period?month=${month}&year=${year}`),
};

export const games = {
  getAll: () => request("/games"),
  getById: (id) => request(`/games/${id}`),
  create: (data) =>
    request("/games", { method: "POST", body: JSON.stringify(data) }),
  update: (id, data) =>
    request(`/games/${id}`, { method: "PUT", body: JSON.stringify(data) }),
  delete: (id) => request(`/games/${id}`, { method: "DELETE" }),
};

export const draws = {
  getAll: () => request("/draws"),
  getById: (id) => request(`/draws/${id}`),
  create: (data) =>
    request("/draws", { method: "POST", body: JSON.stringify(data) }),
  updateResults: (id, data) =>
    request(`/draws/${id}/results`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),
  markAsProcessed: (id) => request(`/draws/${id}/process`, { method: "PUT" }),
  delete: (id) => request(`/draws/${id}`, { method: "DELETE" }),
  getByGame: (id) => request(`/draws/game/${id}`),
  getPending: () => request("/draws/pending"),
};

export const tickets = {
  getAll: () => request("/tickets"),
  getById: (id) => request(`/tickets/${id}`),
  create: (data) =>
    request("/tickets", { method: "POST", body: JSON.stringify(data) }),
  updatePrize: (id, data) =>
    request(`/tickets/${id}/prize`, {
      method: "PUT",
      body: JSON.stringify(data),
    }),
  delete: (id) => request(`/tickets/${id}`, { method: "DELETE" }),
  getByDraw: (id) => request(`/tickets/draw/${id}`),
};
