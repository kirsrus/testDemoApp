import {Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {User} from './models/user';
import {Router} from '@angular/router';
import {AuthenticationService} from './services/authentication.service';
import {Title} from '@angular/platform-browser';
import {SharedService} from './services/shared.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  title = {
    ru: 'МЛТ - медицинские технологии - панель управления',
    en: 'MLT - medical technologies - control panel',
    zh: 'MLT - 医疗技术-控制面板',
  };

  currentThemeLight = true;

  // Текущий язык
  currentLang = 'ru';

  // Индикатор открытой панели навигации
  sidenavOpened = true;

  // Текущий залогиненный пользователь (или null)
  currentUser: User;

  constructor(
    private router: Router,
    private authenticationService: AuthenticationService,
    private elRef: ElementRef,
    private titleService: Title,
    private sharedService: SharedService,
  ) {
    // Привязываем обновление данных о юзере
    this.authenticationService.currentUser.subscribe(u => this.currentUser = u);
  }

  ngOnInit(): void {
    this.sharedService.setTitle(this.title);
    this.currentLang = this.sharedService.currentLang();

    window.document.getElementById('logo').style.backgroundImage = 'url(assets/brand/mlt-logo-' + this.currentLang + '.svg)';

  }

  // Разлогиниваемся и переходим на страницу ввода логина
  logout(): void {
    this.authenticationService.logout();
    this.router.navigate(['/login']);
  }

  // Смена темы
  selectTheme(): void {
    const cn = this.elRef.nativeElement.ownerDocument.body;
    if (cn.className === 'unicorn-dark-theme') {
      cn.className = 'candy-light-theme';
      this.currentThemeLight = true;
      window.document.getElementById('logo').style.backgroundImage = 'url(assets/brand/mlt-logo-' + this.currentLang + '.svg)';
    } else {
      cn.className = 'unicorn-dark-theme';
      this.currentThemeLight = false;
      window.document.getElementById('logo').style.backgroundImage = 'url(assets/brand/mlt-logo-' + this.currentLang + '-dark.svg)';
    }


  }


}
