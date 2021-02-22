export const environment = {
  production: true,
  address: window.location.href.replace(new RegExp('/$', ''), ''),
  addressGraphQL: new URL(window.location.href).origin,
};
