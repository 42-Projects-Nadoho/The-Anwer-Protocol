# World Design

This document provides a comprehensive overview of the game world, including room connections, item distribution, NPCs, and available quests.

## Map Overview

The map consists of an 8-room loop, allowing continuous exploration, with one optional branch leading to the Secret Cave. 

```mermaid
flowchart TD
    classDef room fill:#fff,stroke:#000,stroke-width:2px;

    SC["<div style='text-align: left'><b>Secret Cave</b><hr><i>(Optional Branch)</i></div>"]:::room
    DI["<div style='text-align: left'><b>Destiny Islands</b><hr>+ Items: Potion<br>+ NPCs: Master Yen Sid</div>"]:::room
    TT["<div style='text-align: left'><b>Traverse Town</b><hr>+ Items: Kingdom Key<br>+ NPCs: Leon</div>"]:::room
    WL["<div style='text-align: left'><b>Wonderland</b><hr><i>(Empty)</i></div>"]:::room
    OC["<div style='text-align: left'><b>Olympus Coliseum</b><hr>+ NPCs: Large Body<br>+ NPCs: Sephiroth</div>"]:::room
    A["<div style='text-align: left'><b>Agrabah</b><hr>+ Items: Ether</div>"]:::room
    HT["<div style='text-align: left'><b>Halloween Town</b><hr><i>(Empty)</i></div>"]:::room
    NL["<div style='text-align: left'><b>Neverland</b><hr>+ Items: Wayfinder</div>"]:::room
    HB["<div style='text-align: left'><b>Hollow Bastion</b><hr>+ NPCs: Shadow</div>"]:::room

    SC --South--> DI 

    DI --North--> SC
    DI --East--> TT
    DI --South--> HB

    TT --West--> DI
    TT --East--> WL

    WL --West--> TT
    WL --South--> OC

    OC --North--> WL
    OC --South--> A

    A --North--> OC
    A --West--> HT

    HT --East--> A
    HT --West--> NL

    NL --East--> HT
    NL --North--> HB

    HB --South--> NL
    HB --North--> DI

```

## NPCs

| NPC ID | Name | Role | Location | Description | Stats |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `yen_sid` | Master Yen Sid | Quest Giver | Destiny Islands | A wise and powerful sorcerer, retired Keyblade Master. | - |
| `leon` | Leon | Dialogue | Traverse Town | A stoic warrior wielding a gunblade. | - |
| `shadow_heartless` | Shadow | Enemy | Hollow Bastion | A pureblood Heartless born from the darkness in a heart. | HP: 30, DMG: 10 |
| `large_body` | Large Body | Enemy | Olympus Coliseum | A massive Heartless that blocks attacks with its big belly. | HP: 100, DMG: 25 |
| `sephiroth` | Sephiroth | Enemy | Olympus Coliseum | A legendary SOLDIER with a massive nodachi. | HP: 9999, DMG: 9999 |

## Items

| Item ID | Name | Description | Location |
| :--- | :--- | :--- | :--- |
| `potion` | Potion | A restorative item that heals a small amount of HP. | Destiny Islands |
| `ether` | Ether | An item that restores magic power. | Agrabah |
| `keyblade` | Kingdom Key | A mysterious weapon shaped like a giant key. | Traverse Town |
| `wayfinder` | Wayfinder | A star-shaped charm made of thalassa shells. | Neverland |

## Quests

| Quest ID | Name | Type | Target | Reward | Given By |
| :--- | :--- | :--- | :--- | :--- | :--- |
| `find_wayfinder` | Bonds of Friendship | Fetch | `wayfinder` | You feel your heart grow stronger. | Master Yen Sid |
| `defeat_shadow` | Push Back the Darkness | Defeat | `shadow_heartless` | You gained valuable combat experience. | Master Yen Sid |
