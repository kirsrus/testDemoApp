import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BehaviorSubject, Observable} from 'rxjs';
import {map} from 'rxjs/operators';
import {User} from '../models/user';
import {environment} from '../../environments/environment';

const userIdInLocalStorage = 'currentUser';
// const loginURL = environment.address + '/login';
const loginURL = environment.address + '/login';

@Injectable({
  providedIn: 'root'
})
export class AuthenticationService {
  private currentUserSubject: BehaviorSubject<User>;
  public currentUser: Observable<User>;

  constructor(private http: HttpClient) {
    this.currentUserSubject = new BehaviorSubject<User>(JSON.parse(localStorage.getItem(userIdInLocalStorage)));
    this.currentUser = this.currentUserSubject.asObservable();
  }

  public get currentUserValue(): User {
    return this.currentUserSubject.value;
  }

  // Логинимся на сервере
  login(username: string, password: string): any {
    return this.http.post<any>(loginURL, {username, password})
      .pipe(map(user => {

        // todo: диагностика
        console.log(user);

        // Сохранение полученных данных авторизации в локальном хранилище браузера
        localStorage.setItem(userIdInLocalStorage, JSON.stringify(user));
        this.currentUserSubject.next(user);

        // todo: diag
        console.log(this.currentUserValue);

        return user;
      }));
  }

  // Разлогиниваемся
  logout(): void {
    // Удаляем пользователя из локального хранилища
    localStorage.removeItem(userIdInLocalStorage);
    this.currentUserSubject.next(null);
  }

}
