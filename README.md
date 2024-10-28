# Almacal – der Kalender aus dem AlmaWeb

Mit diesem Programm lässt sich der Kalender aus dem Almaweb im Google-/Apple-Kalender anzeigen.

## Nutzung

Zunächst müssen Sie Ihren Nutzernamen und Ihr Passwort auf der Basis 64 kodieren. Dazu nutzen Sie bitte ein Programm Ihrer Wahl,
beispielsweise <https://base64encode.org>. Sie kodieren Nutzername und Passwort nach dem Schema `Nutzername:Passwort`.

Als nächstes erzeugen Sie Ihre persönliche Kalender-URL. Dazu fügen Sie an der markierten Stelle im folgenden Link einfach Ihren
eben kodierten Code ein: `https://almacal.kleetec.de/?credentials=HIER_EINFÜGEN`

Als Beispiel: Mein Nutzername sei `max.mustermann` und das Passwort `1234`. Dann gebe ich im Kodierer `max.mustermann:1234` ein
und erhalte `bWF4Lm11c3Rlcm1hbm46MTIzNA==`. Meine URL lautet dann also `https://almacal.kleetec.de/?credentials=bWF4Lm11c3Rlcm1hbm46MTIzNA==`

## Google

Um das Programm zu nutzen, gehen Sie auf der Desktop-Website von Google und dann links unten bei `Weitere Kalender` auf das Plus.
Dann drücken Sie auf `Per URL` und geben in das Textfeld Ihre persönliche URL ein.

## Apple

Öffnen Sie Ihre persönliche URL in Safari. Gehen Sie dann auf `Kalender importieren`.
