#!/usr/bin/env python
# -*- coding: utf-8 -*- 

"""Usage: import.py

Download TV schedules from TheTVDB and output a comma separated file
"""

from __future__ import print_function

import codecs
import json
import sys
import tvdb_api

from datetime import date
from datetime import datetime
from datetime import timedelta

SHOWS = [
    'The Blacklist',
    'Castle',
    'How I Met Your Mother',
    'Suits',
    'The Mentalist',
    'White Collar'
]

def main():
    t = tvdb_api.Tvdb()
    episodes = []
    cutoff = date.today() - timedelta(14)
    for name in SHOWS:
        for x in find_episodes(t, name, cutoff):
            episodes.append(x)
    json.dump(episodes, sys.stdout)


def find_episodes(tvdb, series, cutoff_date):
    """
    Yield episodes aired after the specified cutoff date in reverse order.
    i.e. most recent first.
    """
    try:
        show = tvdb[series]
        last_season = max(show.keys())
        for ep in episodes(show[last_season]):
            dt = airdate(ep)
            if dt and dt >= cutoff_date:
                yield {
                    'series':  series,
                    'season':  last_season,
                    'title':   ep['episodename'],
                    'number':  int(ep['episodenumber']),
                    'airdate': ep['firstaired'] + 'T00:00:00Z'
                }
            else:
                break
    except tvdb_api.tvdb_shownotfound:
        print("{}: Show not found".format(series), file=sys.stderr)


def last_season(show):
    """Returns the most recently aired season of a show."""
    return show[max(show.keys())]


def episodes(season):
    """
    Yield episodes in reverse order of date aired. i.e most recent first
    """
    for i in sorted(list(season.keys()), reverse=True):
        yield season[i]


def airdate(episode):
    try:
        if episode['firstaired']:
            return datetime.strptime(episode['firstaired'], '%Y-%m-%d').date()
        else:
            return None
    except tvdb_api.tvdb_attributenotfound:
        return None


if __name__ == '__main__':
    status = main()
    sys.exit(status)
