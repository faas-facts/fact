/*
 *  MIT License
 *
 *  Copyright (c) 2021. Fact Contributors
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 */
package fact;

import java.util.HashMap;
import java.util.Map;
import com.amazonaws.services.lambda.runtime.Context;
import io.github.fact.Fact;
import io.github.fact.FactConfiguration;
import io.github.fact.FactConfigurationBuilder;

import static java.lang.Thread.sleep;


public class ConsoleFactTest {

    public String handleRequest(Map<String, String> event, Context context) {
        FactConfigurationBuilder builder= new FactConfigurationBuilder();
        FactConfiguration conf= builder.createLazyLogger();
        Fact.boot(conf);
        Fact.start(context,event);
        try {
            sleep(1000);
        } catch (InterruptedException e) { }
        Map<String,String>tags= new HashMap<String, String>();
        tags.put("experiment","#1");
        Fact.update(context,"update testing",tags);
        try {
            sleep(1000);
        } catch (InterruptedException e) { }
        Fact.done(context,"im done with this","no more args");
        return null;
    }
}
