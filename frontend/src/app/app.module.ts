import {NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {MaterialModule} from './material.module';
import {LoginComponent} from './login/login.component';
import {FormsModule} from '@angular/forms';
import {AboutComponent} from './about/about.component';
import {HomeComponent} from './home/home.component';
import {HttpClientModule} from '@angular/common/http';
import {GraphQLModule} from './graphql.module';
import {FlexLayoutModule} from '@angular/flex-layout';
import { HumanizebytesPipe } from './pipes/humanizebytes.pipe';
import { ConsoleComponent } from './console/console.component';


@NgModule({
  declarations: [
    AboutComponent,
    AppComponent,
    LoginComponent,
    HomeComponent,
    HumanizebytesPipe,
    ConsoleComponent,
  ],
  imports: [
    AppRoutingModule,
    BrowserAnimationsModule,
    BrowserModule,
    FormsModule,
    HttpClientModule,
    MaterialModule,
    GraphQLModule,
    FlexLayoutModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {

}
