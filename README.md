# ClinicHub – portal za klinike baziran na mikroservisnoj arhitekturi

Matija Petrović SW 33-2017

## Opis aplikacije

ClinicHub je portal koji omogućuje korisnicima da pretražuju klinike i doktore zaposlene u njima, zakazuju preglede i ocenjuju klinike i doktore nakon obavljenog pregleda. Uloge koje postoje u sistemu su prethodno pomenuti regularan korisnik i administrator kome je omogućeno da ažurira podatke o klinikama, doktorima i pregledima, kao i da odobrava registracije korisnika.

## Arhitektura sistema

Arhitektura koja bi se koristila za implementaciju sistema bi bila zasnovana na mikroservisima za čiji razvoj bi se koristili jezici Go, Pharo i Python. Svaki od mikroservisa bi bio povezan sa svojom MySQL bazom podataka, a Go i Python mikroservisi bi bili implementirani kao REST API. Servis za autentifikaciju bi prilikom logina generisao JWT token koji bi klijenti koristili za pristup ostalim mikroservisima.

## Opis mikroservisa

### User microservice

Python Django/Flask REST servis koji omogućava registraciju i prijavu korisnika na sistem. Nakon prijave generiše JWT token koji se koristi za pristup ostalim servisima. Administratori koriste ovaj mikroservis kako bi odobrili registraciju korisnika.

### Clinic microservice

Go REST servis koji omogućava administratorima da dodaju i ažuriraju klinike, cene pregleda u klinici i doktore koji rade u tim klinikama. Svaki doktor je specijalizovan za određeni tip pregleda i ima određeno radno vreme. Regularni korisnici mogu da pretražuju klinike i doktore na osnovu tipa pregleda koji žele da obave i datuma za koji na klinikama postoje doktori sa slobodnim terminima.

### Scheduling microservice

Go REST servis koji omogućava korisnicima da nakon pretrage i odabira klinike i doktora zakažu pregled za određeni termin. Takođe omogućava _Clinic microservice_-u da proveri koji doktori i klinike odgovaraju zahtevima pretrage korisnika odnosno da li su slobodni za traženi datum.

### Rating microservice

Go REST servis koji omogućava korisnicima da nakon završenog pregleda ocene kliniku i doktora na kojoj su se pregledali.

### Analytics client

Pharo klijent koji omogućava praćenje poslovanja klinika, prikazuje prosek prihoda i ocena za kliniku u traženom vremenskom periodu.

### Vue client

Vue.js klijent koji omogućava pristup svim funkcionalnostima sistema.

## Dodatne ideje za diplomski

- API gateway za uniforman pristup mikroservisima
- OAuth2 protokol za autentifikaciju i autorizaciju
- Nakon zakazivanja, administrator odobrava pregled i dodeljuje mu jednu od slobodnih soba u okviru klinike + mikroservis koji na kraju dana automatski dodeljuje sobe neodobrenim zahtevima za pregled
- Kontejnerizacija mikroservisa
