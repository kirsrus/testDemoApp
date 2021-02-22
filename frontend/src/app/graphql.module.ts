import {NgModule} from '@angular/core';
import {APOLLO_OPTIONS} from 'apollo-angular';
import {ApolloClientOptions, InMemoryCache} from '@apollo/client/core';
import {HttpLink} from 'apollo-angular/http';
import {environment} from '../environments/environment';
import {WebSocketLink} from '@apollo/client/link/ws';


// const cURL = new URL(window.location.href);
// const uri = 'ws://' + cURL.origin + '/graphql'; // <-- add the URL of the GraphQL server here
const uri = (environment.addressGraphQL + '/graphql').replace('http://', 'ws://');
console.log('GraphQL URI:', uri);


// export function createApollo(httpLink: HttpLink): ApolloClientOptions<any> {
//   return {
//     link: httpLink.create({uri}),
//     cache: new InMemoryCache(),
//   };
// }

export function createApollo(): ApolloClientOptions<any> {
  return {
    cache: new InMemoryCache(),
    link: new WebSocketLink({
      uri,
      options: {
        reconnect: true,
        connectionParams: {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('token')}`
          }
        }
      }
    }),
  };
}

@NgModule({
  providers: [
    {
      provide: APOLLO_OPTIONS,
      useFactory: createApollo,
      deps: [HttpLink],
    },
  ],
})
export class GraphQLModule {
}
