# Lineschedule

Lineschedule is a crude tool to generate custom transit timetables from static [GTFS data](https://developers.google.com/transit/gtfs) published by some public transit authorities.

Between the pandemic and the neoliberal asset stripping of public services in most western countries, Lineschedule has few practical uses beyond documenting the decline of public transit services.

In addition, this software is very likely to contain bugs. Don't use it to plan trips that are important.

```
> lineschedule --help
lineschedule

  Usage:
    lineschedule [import|information|schedule|csvschedule]

  Subcommands: 
    import        Import GTFS data from directory into database.
    information   Get information for route.
    schedule      Get schedule for route.
    csvschedule   Make line schedule CSV(s) from XML configuration.

  Flags: 
       --version   Displays the program version string.
    -h --help      Displays help with available flag, subcommand, and positional value parameters.
    -db --dbname    Override database filename (default: gtfs.sqlite3)
```

## Example workflow

```
> lineschedule --import GTFSFeed
```

Imports a GTFS dataset from the directory GTFSFeed into Lineschedule's working database. This is a required first-step before using any other Lineschedule functions.

Further examples use the GTFS feed from [Translink (BC)](https://translink.ca), which operates public transit in the Metro Vancouver area.

```
> lineschedule information -r R2 -s 2022-06-02 
Service on route R2 for 2022-06-02:

Direction 0:
Stop ID  Trips at stop  Stop name
4461     95/95          Park Royal @ Bay 5
4463     95/95          Eastbound Marine Dr @ Capilano Rd
12685    95/95          Eastbound Marine Dr @ Pemberton Ave
4456     95/95          Eastbound Marine Dr @ Hamilton Ave
11929    95/95          Eastbound W 3rd St @ Bewicke Ave
4459     95/95          Lonsdale Quay @ Bay 7
4212     95/95          Eastbound E 3rd St @ Lonsdale Ave
3994     95/95          Eastbound E 3rd St @ Ridgeway Ave
11598    95/95          Eastbound Cotton Rd @ Brooksbank Ave
4143     95/95          Phibbs Exchange @ Bay 2

Direction 1:
Stop ID  Trips at stop  Stop name
4143     95/95          Phibbs Exchange @ Bay 2
11597    95/95          Westbound Cotton Rd @ Brooksbank Ave
4194     95/95          Westbound E 3rd St @ Ridgeway Ave
12682    95/95          Southbound Lonsdale Ave @ W 3rd St
10865    95/95          Lonsdale Quay @ Bay 3
4427     95/95          Westbound Marine Dr @ Bewicke Ave
11367    95/95          Westbound Marine Dr @ Hamilton Ave
4485     95/95          Westbound Marine Dr @ Pemberton Ave
12683    95/95          Westbound Marine Dr @ Capilano Rd
11866    95/95          Park Royal @ Bay 1
11868    95/95          Eastbound Marine Dr @ South Mall Access
```
The `information` subcommand displays all stops served by a route (specified by `-r`), in both directions, for the date specified (by `-s`).

```
> lineschedule schedule -r R2 -s 2022-06-02 -d 1 -st 4143 -st 12683 
              4143: Phibbs Exchange @ Bay 2 12683: Westbound Marine Dr @ Capilano Rd 
R2 PARK ROYAL 05:10                         05:32                                    
R2 PARK ROYAL 05:25                         05:47                                    
R2 PARK ROYAL 05:37                         05:59                                    
R2 PARK ROYAL 05:47                         06:09                                    
R2 PARK ROYAL 05:57                         06:19                                    
R2 PARK ROYAL 06:07                         06:30                                    
R2 PARK ROYAL 06:17                         06:40                                    
[... Omitted ...]

```

The `schedule` subcommand displays a timetable for the stops specified by `-st` (or all stops if `-st` is omitted), served by a route (specified by `-r`, for the direction specified by `-d`, for the date specified by `-s`.

## Scripted Workflow

Lineschedule's major function is to generate complex timetables, including trips with connections, using an XML configuration file. These timetables are written in .CSV format, for import into a spreadsheet or other downstream tooling.

```
> lineschedule csvschedule -if example.xml
```

Example.xml:
```
<LineSchedule>
  <Journey>
      <Output>Outbound.csv</Output>
      <!-- ServiceDate must be in YYYY-MM-DD format. -->
      <ServiceDate>2022-06-02</ServiceDate>
      <Route>
          <Route>R2</Route>
          <Direction>1</Direction>
          <Stops>10865</Stops>
          <Stops>11367</Stops>
          <Stops>11866</Stops>
      </Route>
      <Route>
          <Route>257</Route>
          <Direction>1</Direction>
          <Stops>4608</Stops>
          <Stops>9129</Stops>
          <Stops>4491</Stops>
      </Route>
      <!-- Defines a connection between two routes that don't share a common stop. -->
      <Connection>
          <FromRoute>R2</FromRoute>
          <FromStop>11866</FromStop>
          <ToRoute>257</ToRoute>
          <ToStop>4491</ToStop>
      </Connection>
      <!-- StopName blocks override stop names given in the GTFS data. -->
      <StopName>
        <StopID>9129</StopID>
        <Name>Seymour</Name>
      </StopName>
      <StopName>
        <StopID>10865</StopID>
        <Name>Lonsdale</Name>
      </StopName>
      <StopName>
        <StopID>11367</StopID>
        <Name>Cap Mall</Name>
      </StopName>
      <StopName>
        <StopID>4608</StopID>
        <Name>Horseshoe Bay</Name>
      </StopName>
  </Journey>
  <Journey>
      <Output>Inbound.csv</Output>
      <ServiceDate>2022-06-02</ServiceDate>
      <Route>
          <Route>257</Route>
          <Direction>0</Direction>
          <Stops>11548</Stops>
          <Stops>4608</Stops>
          <Stops>4461</Stops>
          <Stops>985</Stops>
      </Route>
      <Route>
          <Route>R2</Route>
          <Direction>0</Direction>
          <Stops>4461</Stops>
          <Stops>4456</Stops>
          <Stops>4459</Stops>
      </Route> 
      <StopName>
        <StopID>4608</StopID>
        <Name>Horseshoe Bay</Name>
      </StopName>
      <StopName>
        <StopID>4461</StopID>
        <Name>Park Royal</Name>
      </StopName>
      <StopName>
        <StopID>4456</StopID>
        <Name>Cap Mall</Name>
      </StopName>
      <StopName>
        <StopID>4459</StopID>
        <Name>Lonsdale Quay</Name>
      </StopName>
      <StopName>
        <StopID>985</StopID>
        <Name>Granville</Name>
      </StopName>
  </Journey>
</LineSchedule>
```
## Known Bugs

Lineschedule is very crude and will emit a stack trace, or nonsensical data, in the event of most errors. Due to the limited practical utility of this software, it is not worth the time and effort to expand the codebase to emit human-friendly error messages.

## License

Copyright 2022 Coridon Henshaw

Permission is granted to all natural persons to execute, distribute, and/or modify this software (including its documentation) subject to the following terms:

1. Subject to point \#2, below, **all commercial use and distribution is prohibited.** This software has been released for personal and academic use for the betterment of society through any purpose that does not create income or revenue. *It has not been made available for businesses to profit from unpaid labor.*

2. Re-distribution of this software on for-profit, public use, repository hosting sites (for example: Github) is permitted provided no fees are charged specifically to access this software.

3. **This software is provided on an as-is basis and may only be used at your own risk.** This software is the product of a single individual's recreational project. The author does not have the resources to perform the degree of code review, testing, or other verification required to extend any assurances that this software is suitable for any purpose, or to offer any assurances that it is safe to execute without causing data loss or other damage.

4. **This software is intended for experimental use in situations where data loss (or any other undesired behavior) will not cause unacceptable harm.** Users with critical data safety needs must not use this software and, instead, should use equivalent tools that have a proven track record.

5. If this software is redistributed, this copyright notice and license text must be included without modification.

6. Distribution of modified copies of this software is discouraged but is not prohibited. It is strongly encouraged that fixes, modifications, and additions be submitted for inclusion into the main release rather than distributed independently.

7. This software reverts to the public domain 10 years after its final update or immediately upon the death of its author, whichever happens first.
