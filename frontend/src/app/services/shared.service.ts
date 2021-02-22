import {Injectable} from '@angular/core';
import {Title} from '@angular/platform-browser';

// Описание заголовков страницы
export interface PageTitle {
  ru: string;
  en: string;
  zh: string;
}

@Injectable({
  providedIn: 'root'
})
export class SharedService {

  constructor(private titleService: Title) {
  }

  // Возвращает активный язык из текущего URL
  currentLang(): string {
    const url = new URL(window.location.href);
    if (url.pathname.indexOf('/ru/') === 0) {
      return 'ru';
    } else if (url.pathname.indexOf('/en/') === 0) {
      return 'en';
    } else if (url.pathname.indexOf('/zh/') === 0) {
      return 'zh';
    } else {
      return 'ru';
    }
  }

  setTitle(title: PageTitle): void {
    this.titleService.setTitle(title[this.currentLang()]);
  }

}
