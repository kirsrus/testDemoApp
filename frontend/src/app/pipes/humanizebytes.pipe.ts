import {Pipe, PipeTransform} from '@angular/core';


@Pipe({
  name: 'humanizebytes'
})
export class HumanizebytesPipe implements PipeTransform {

  // transform(value: unknown, ...args: unknown[]): unknown {
  transform(value: number): string {
    return this.number_size_to_human(value);
  }

  /**
   * Перевод текстового представления раземера в цифровой (формата "100 KB")
   * @param strSize размер в виде текста
   * @return размер в цифрах
   */
  human_size_to_number(strSize: string): number {
    let result = 0;
    const m = strSize.match(/^\s*([\d.]+)\s*(\w+)/i);
    if (m) {
      result = parseFloat(m[1]);
      switch (m[2].toUpperCase()) {
        case 'B':
          break;
        case 'KB':
          result = result * 1e+3;
          break;
        case 'MB':
          result = result * 1e+6;
          break;
        case 'GB':
          result = result * 1e+9;
          break;
        case 'TB':
          result = result * 1e+12;
          break;
        case 'PB':
          result = result * 1e+15;
          break;
        case 'EB':
          result = result * 1e+18;
          break;
        case 'ZB':
          result = result * 1e+21;
          break;
        case 'YB':
          result = result * 1e+24;
          break;
      }
    }
    return result;
  }

  /**
   * Переовд числового представления размера в человекопонятный
   * @param size размер в виде числа
   * @return человекопонятное представление
   */
  number_size_to_human(size: number): string {
    const fileSizeName = [' B', ' KB', ' MB', ' GB', ' TB', ' PB', ' EB', ' ZB', ' YB'];
    const logB = (num: number, base: number) => {
      return Math.log(num) / Math.log(base);
    };
    let result: string;
    if (size) {
      const i = Math.floor(logB(size, 1000));  // 1024 переделано на 1000
      const sizeRes = size / Math.pow(1000, (i));  // 1024 переделано на 1000
      result = '' + parseFloat(parseFloat(sizeRes.toString().replace(/e.+/i, '')).toFixed(2)) + fileSizeName[i];
    } else {
      result = '0 B';
    }
    return result;
  }

}
